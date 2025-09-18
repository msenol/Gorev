// Template Wizard JavaScript

(function() {
    const vscode = acquireVsCodeApi();
    
    // State
    let currentStep = 'template-selection';
    let selectedTemplate = null;
    let templates = [];
    let formValues = {};
    
    // Elements
    const elements = {
        // Steps
        stepTemplateSelection: document.getElementById('step-template-selection'),
        stepFormFields: document.getElementById('step-form-fields'),
        stepPreview: document.getElementById('step-preview'),
        
        // Template selection
        templateSearch: document.getElementById('template-search'),
        templateGrid: document.getElementById('template-grid'),
        categoryTabs: document.querySelectorAll('.category-tab'),
        
        // Form
        templateName: document.getElementById('template-name'),
        templateDescription: document.getElementById('template-description'),
        formFields: document.getElementById('form-fields'),
        
        // Preview
        taskPreview: document.getElementById('task-preview'),
        
        // Actions
        btnBack: document.getElementById('btn-back'),
        btnPreview: document.getElementById('btn-preview'),
        btnCreate: document.getElementById('btn-create'),
        btnBackToForm: document.getElementById('btn-back-to-form'),
        btnConfirmCreate: document.getElementById('btn-confirm-create')
    };
    
    // Initialize
    function init() {
        // Attach event listeners
        elements.templateSearch.addEventListener('input', debounce(handleSearch, 300));
        
        elements.categoryTabs.forEach(tab => {
            tab.addEventListener('click', handleCategoryClick);
        });
        
        elements.btnBack.addEventListener('click', () => showStep('template-selection'));
        elements.btnPreview.addEventListener('click', handlePreview);
        elements.btnCreate.addEventListener('click', handleCreate);
        elements.btnBackToForm.addEventListener('click', () => showStep('form-fields'));
        elements.btnConfirmCreate.addEventListener('click', handleCreate);
        
        // Load initial templates
        vscode.postMessage({ command: 'loadTemplates' });
    }
    
    // Step navigation
    function showStep(step) {
        currentStep = step;
        
        // Hide all steps
        document.querySelectorAll('.wizard-step').forEach(s => {
            s.classList.remove('active');
        });
        
        // Show current step
        switch (step) {
            case 'template-selection':
                elements.stepTemplateSelection.classList.add('active');
                break;
            case 'form-fields':
                elements.stepFormFields.classList.add('active');
                break;
            case 'preview':
                elements.stepPreview.classList.add('active');
                break;
        }
    }
    
    // Template selection handlers
    function handleSearch(event) {
        const query = event.target.value.trim();
        if (query) {
            vscode.postMessage({ command: 'searchTemplates', query });
        } else {
            const activeCategory = document.querySelector('.category-tab.active').dataset.category;
            loadTemplatesByCategory(activeCategory);
        }
    }
    
    function handleCategoryClick(event) {
        // Update active tab
        elements.categoryTabs.forEach(tab => tab.classList.remove('active'));
        event.target.classList.add('active');

        const category = event.target.dataset.category;

        if (category === 'favorites') {
            loadFavoriteTemplates();
        } else {
            loadTemplatesByCategory(category);
        }
    }

    function loadFavoriteTemplates() {
        // Show loading state
        showLoadingState(elements.templateGrid, 'Favori şablonlar yükleniyor...');

        // Get all templates first, then filter favorites
        vscode.postMessage({ command: 'loadTemplates' });
    }
    
    function loadTemplatesByCategory(category) {
        // Show loading state
        showLoadingState(elements.templateGrid, 'Şablonlar yükleniyor...');

        vscode.postMessage({
            command: 'loadTemplates',
            category: category || undefined
        });
    }

    function showLoadingState(container, message = 'Yükleniyor...') {
        container.innerHTML = `
            <div class="loading-state">
                <div class="loading-spinner"></div>
                <div class="loading-message">${message}</div>
            </div>
        `;
        container.classList.add('loading');
    }

    function showErrorState(container, message, retry = null) {
        container.innerHTML = `
            <div class="error-state">
                <div class="error-icon">⚠️</div>
                <div class="error-message">${message}</div>
                ${retry ? `<button class="retry-btn" onclick="${retry}">Tekrar Dene</button>` : ''}
            </div>
        `;
        container.classList.remove('loading');
        container.classList.add('error');
    }
    
    function renderTemplates(templateList) {
        templates = templateList;

        // Remove loading state
        elements.templateGrid.classList.remove('loading');

        if (templates.length === 0) {
            elements.templateGrid.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">📋</div>
                    <div class="empty-state-message">Şablon bulunamadı</div>
                    <div class="empty-state-hint">Farklı bir kategori veya arama terimi deneyin</div>
                </div>
            `;
            return;
        }
        
        elements.templateGrid.innerHTML = templates.map(template => `
            <div class="template-card" data-template-id="${template.id}">
                <button class="template-favorite ${isFavorite(template.id) ? 'active' : ''}" 
                        onclick="toggleFavorite('${template.id}', event)">
                    ${isFavorite(template.id) ? '⭐' : '☆'}
                </button>
                <div class="template-icon">${getTemplateIcon(template.kategori)}</div>
                <div class="template-name">${template.isim}</div>
                <div class="template-category">${template.kategori || 'Genel'}</div>
                <div class="template-description">${template.tanim || ''}</div>
            </div>
        `).join('');
        
        // Attach click handlers
        document.querySelectorAll('.template-card').forEach(card => {
            card.addEventListener('click', handleTemplateSelect);
        });
    }
    
    function handleTemplateSelect(event) {
        // Ignore if clicking on favorite button
        if (event.target.classList.contains('template-favorite')) {
            return;
        }
        
        const templateId = event.currentTarget.dataset.templateId;
        vscode.postMessage({ command: 'selectTemplate', templateId });
    }
    
    function getTemplateIcon(category) {
        const icons = {
            'Genel': '📝',
            'Teknik': '🔧',
            'Özellik': '✨',
            'Bug': '🐛',
            'Araştırma': '🔍',
            'Dokümantasyon': '📚'
        };
        return icons[category] || '📋';
    }

    // Enhanced field renderers
    function renderTextField(field, fieldId, required, value) {
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <input type="text" id="${fieldId}" name="${field.isim}"
                       class="form-input enhanced-input" value="${escapeHtml(value)}" ${required}
                       placeholder="${field.placeholder || ''}"
                       autocomplete="off">
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Bu alan zorunludur</div>
                <div class="validation-feedback"></div>
            </div>
        `;
    }

    function renderTextareaField(field, fieldId, required, value) {
        const isLarge = field.buyuk || field.isim.toLowerCase().includes('aciklama');
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <textarea id="${fieldId}" name="${field.isim}"
                          class="form-textarea enhanced-textarea ${isLarge ? 'large' : ''}" ${required}
                          placeholder="${field.placeholder || 'Detayları buraya yazın...'}"
                          rows="${isLarge ? 8 : 4}">${escapeHtml(value)}</textarea>
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Bu alan zorunludur</div>
                <div class="character-count"><span class="count">0</span> karakter</div>
            </div>
        `;
    }

    function renderSelectField(field, fieldId, required, value) {
        const options = field.secenekler || ['dusuk', 'orta', 'yuksek'];
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <div class="select-wrapper">
                    <select id="${fieldId}" name="${field.isim}"
                            class="form-select enhanced-select" ${required}>
                        <option value="">Seçim yapın...</option>
                        ${options.map(opt =>
                            `<option value="${escapeHtml(opt)}" ${value === opt ? 'selected' : ''}>${opt}</option>`
                        ).join('')}
                    </select>
                    <div class="select-arrow">▼</div>
                </div>
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Bu alan zorunludur</div>
            </div>
        `;
    }

    function renderDateField(field, fieldId, required, value) {
        const today = new Date().toISOString().split('T')[0];
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <div class="date-input-wrapper">
                    <input type="date" id="${fieldId}" name="${field.isim}"
                           class="form-input date-input" value="${value}" ${required}
                           min="${today}">
                    <div class="date-shortcuts">
                        <button type="button" class="date-shortcut" data-days="1">Yarın</button>
                        <button type="button" class="date-shortcut" data-days="7">1 Hafta</button>
                        <button type="button" class="date-shortcut" data-days="30">1 Ay</button>
                    </div>
                </div>
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Bu alan zorunludur</div>
            </div>
        `;
    }

    function renderTagsField(field, fieldId, required, value) {
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <div class="tag-input-container enhanced-tags" id="${fieldId}-container">
                    <input type="text" id="${fieldId}" class="tag-input"
                           placeholder="Etiket ekle ve Enter'a bas...">
                </div>
                <input type="hidden" name="${field.isim}" value="${escapeHtml(value)}">
                <div class="tag-suggestions">
                    <div class="common-tags">
                        <span class="tag-suggestion" data-tag="bug">bug</span>
                        <span class="tag-suggestion" data-tag="feature">feature</span>
                        <span class="tag-suggestion" data-tag="urgent">urgent</span>
                        <span class="tag-suggestion" data-tag="documentation">documentation</span>
                    </div>
                </div>
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Bu alan zorunludur</div>
            </div>
        `;
    }

    function renderEmailField(field, fieldId, required, value) {
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <input type="email" id="${fieldId}" name="${field.isim}"
                       class="form-input enhanced-input" value="${escapeHtml(value)}" ${required}
                       placeholder="ornek@email.com">
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Geçerli bir email adresi girin</div>
                <div class="validation-feedback"></div>
            </div>
        `;
    }

    function renderUrlField(field, fieldId, required, value) {
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <input type="url" id="${fieldId}" name="${field.isim}"
                       class="form-input enhanced-input" value="${escapeHtml(value)}" ${required}
                       placeholder="https://example.com">
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Geçerli bir URL girin</div>
                <div class="validation-feedback"></div>
            </div>
        `;
    }

    function renderNumberField(field, fieldId, required, value) {
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <input type="number" id="${fieldId}" name="${field.isim}"
                       class="form-input enhanced-input" value="${value}" ${required}
                       min="${field.min || 0}" max="${field.max || ''}"
                       step="${field.step || 1}">
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Geçerli bir sayı girin</div>
                <div class="validation-feedback"></div>
            </div>
        `;
    }

    function renderMarkdownField(field, fieldId, required, value) {
        return `
            <div class="form-group">
                <label class="form-label ${required}" for="${fieldId}">
                    ${field.isim}
                    ${field.zorunlu ? '<span class="required-indicator">*</span>' : ''}
                </label>
                <div class="markdown-editor">
                    <div class="markdown-toolbar">
                        <button type="button" class="md-btn" data-action="bold" title="Kalın"><b>B</b></button>
                        <button type="button" class="md-btn" data-action="italic" title="İtalik"><i>I</i></button>
                        <button type="button" class="md-btn" data-action="link" title="Link">🔗</button>
                        <button type="button" class="md-btn" data-action="list" title="Liste">📝</button>
                        <div class="md-divider"></div>
                        <button type="button" class="md-btn preview-toggle" data-target="${fieldId}">👁 Önizleme</button>
                    </div>
                    <textarea id="${fieldId}" name="${field.isim}"
                              class="form-textarea markdown-textarea" ${required}
                              placeholder="Markdown formatında yazın..."
                              rows="8">${escapeHtml(value)}</textarea>
                    <div class="markdown-preview" id="${fieldId}-preview" style="display: none;"></div>
                </div>
                ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                <div class="form-error">Bu alan zorunludur</div>
            </div>
        `;
    }

    // Utility function for HTML escaping
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // Form handling
    function renderForm(template) {
        selectedTemplate = template;

        // Update header
        elements.templateName.textContent = template.isim;
        elements.templateDescription.textContent = template.tanim || '';

        // Build form fields with enhanced rendering
        const fieldsHtml = template.alanlar.map(field => {
            const fieldId = `field-${field.isim}`;
            const required = field.zorunlu ? 'required' : '';
            const value = formValues[field.isim] || field.varsayilan || '';

            switch (field.tur) {
                case 'text':
                    return renderTextField(field, fieldId, required, value);

                case 'textarea':
                    return renderTextareaField(field, fieldId, required, value);

                case 'select':
                    return renderSelectField(field, fieldId, required, value);

                case 'date':
                    return renderDateField(field, fieldId, required, value);

                case 'tags':
                    return renderTagsField(field, fieldId, required, value);

                case 'email':
                    return renderEmailField(field, fieldId, required, value);

                case 'url':
                    return renderUrlField(field, fieldId, required, value);

                case 'number':
                    return renderNumberField(field, fieldId, required, value);

                case 'markdown':
                    return renderMarkdownField(field, fieldId, required, value);

                default:
                    return renderTextField(field, fieldId, required, value);
            }
        }).join('');
        
        elements.formFields.innerHTML = fieldsHtml;
        
        // Initialize enhanced form features
        initializeEnhancedForm();

        // Show form step
        showStep('form-fields');
    }

    function initializeEnhancedForm() {
        // Initialize tag inputs
        document.querySelectorAll('.tag-input').forEach(input => {
            initializeTagInput(input);
        });

        // Initialize character counters
        document.querySelectorAll('.enhanced-textarea').forEach(textarea => {
            initializeCharacterCounter(textarea);
        });

        // Initialize date shortcuts
        document.querySelectorAll('.date-shortcut').forEach(btn => {
            btn.addEventListener('click', handleDateShortcut);
        });

        // Initialize tag suggestions
        document.querySelectorAll('.tag-suggestion').forEach(suggestion => {
            suggestion.addEventListener('click', handleTagSuggestion);
        });

        // Initialize markdown editors
        document.querySelectorAll('.markdown-editor').forEach(editor => {
            initializeMarkdownEditor(editor);
        });

        // Initialize real-time validation
        document.querySelectorAll('.enhanced-input, .enhanced-textarea, .enhanced-select').forEach(field => {
            field.addEventListener('input', handleRealtimeValidation);
            field.addEventListener('blur', handleFieldBlur);
        });
    }

    function initializeCharacterCounter(textarea) {
        const counter = textarea.parentElement.querySelector('.character-count .count');
        if (counter) {
            const updateCounter = () => {
                counter.textContent = textarea.value.length;
            };
            textarea.addEventListener('input', updateCounter);
            updateCounter();
        }
    }

    function handleDateShortcut(event) {
        const days = parseInt(event.target.dataset.days);
        const dateInput = event.target.closest('.date-input-wrapper').querySelector('.date-input');
        const futureDate = new Date();
        futureDate.setDate(futureDate.getDate() + days);
        dateInput.value = futureDate.toISOString().split('T')[0];

        // Trigger validation
        dateInput.dispatchEvent(new Event('input', { bubbles: true }));
    }

    function handleTagSuggestion(event) {
        const tag = event.target.dataset.tag;
        const container = event.target.closest('.form-group').querySelector('.tag-input-container');
        const tagInput = container.querySelector('.tag-input');

        addTag(container, tag);
        updateTagsInput(container);

        // Remove suggestion after use
        event.target.style.opacity = '0.5';
        event.target.style.pointerEvents = 'none';
    }

    function initializeMarkdownEditor(editor) {
        const toolbar = editor.querySelector('.markdown-toolbar');
        const textarea = editor.querySelector('.markdown-textarea');
        const preview = editor.querySelector('.markdown-preview');

        // Toolbar actions
        toolbar.addEventListener('click', (event) => {
            const action = event.target.dataset.action;
            if (action) {
                handleMarkdownAction(textarea, action);
            }
        });

        // Preview toggle
        const previewToggle = toolbar.querySelector('.preview-toggle');
        if (previewToggle) {
            previewToggle.addEventListener('click', () => {
                if (preview.style.display === 'none') {
                    // Show preview
                    if (typeof marked !== 'undefined') {
                        preview.innerHTML = marked.parse(textarea.value);
                    } else {
                        preview.innerHTML = '<p>Markdown önizlemesi yüklenemedi</p>';
                    }
                    preview.style.display = 'block';
                    textarea.style.display = 'none';
                    previewToggle.textContent = '✏️ Düzenle';
                } else {
                    // Show editor
                    preview.style.display = 'none';
                    textarea.style.display = 'block';
                    previewToggle.textContent = '👁 Önizleme';
                }
            });
        }
    }

    function handleMarkdownAction(textarea, action) {
        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const selectedText = textarea.value.substring(start, end);
        let replacement = '';

        switch (action) {
            case 'bold':
                replacement = `**${selectedText || 'kalın metin'}**`;
                break;
            case 'italic':
                replacement = `*${selectedText || 'italik metin'}*`;
                break;
            case 'link':
                replacement = `[${selectedText || 'link metni'}](https://example.com)`;
                break;
            case 'list':
                replacement = selectedText
                    ? selectedText.split('\n').map(line => `- ${line}`).join('\n')
                    : '- Liste öğesi';
                break;
        }

        if (replacement) {
            textarea.setRangeText(replacement, start, end, 'end');
            textarea.focus();
        }
    }

    function handleRealtimeValidation(event) {
        const field = event.target;
        const formGroup = field.closest('.form-group');
        const errorElement = formGroup.querySelector('.form-error');
        const feedbackElement = formGroup.querySelector('.validation-feedback');

        // Clear previous validation state
        formGroup.classList.remove('error', 'success');
        if (feedbackElement) feedbackElement.textContent = '';

        // Validate based on field type
        const isValid = validateField(field);

        if (field.hasAttribute('required') && !field.value.trim()) {
            // Required field is empty - don't show error yet if user is still typing
            return;
        }

        if (field.value.trim() && !isValid) {
            formGroup.classList.add('error');
            if (feedbackElement) {
                feedbackElement.textContent = getValidationMessage(field);
            }
        } else if (field.value.trim() && isValid) {
            formGroup.classList.add('success');
            if (feedbackElement) {
                feedbackElement.textContent = '✓';
            }
        }
    }

    function handleFieldBlur(event) {
        const field = event.target;
        const formGroup = field.closest('.form-group');

        // Show required field error on blur if empty
        if (field.hasAttribute('required') && !field.value.trim()) {
            formGroup.classList.add('error');
        }
    }

    function validateField(field) {
        const value = field.value.trim();

        switch (field.type) {
            case 'email':
                return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
            case 'url':
                try {
                    new URL(value);
                    return true;
                } catch {
                    return false;
                }
            case 'number':
                const num = parseFloat(value);
                const min = field.min ? parseFloat(field.min) : -Infinity;
                const max = field.max ? parseFloat(field.max) : Infinity;
                return !isNaN(num) && num >= min && num <= max;
            default:
                return true;
        }
    }

    function getValidationMessage(field) {
        switch (field.type) {
            case 'email':
                return 'Geçerli bir email adresi girin';
            case 'url':
                return 'Geçerli bir URL girin (http:// veya https:// ile başlamalı)';
            case 'number':
                return 'Geçerli bir sayı girin';
            default:
                return 'Bu alan geçerli değil';
        }
    }
    
    function initializeTagInput(input) {
        const container = input.parentElement;
        const hiddenInput = container.parentElement.querySelector('input[type="hidden"]');
        const tags = hiddenInput.value ? hiddenInput.value.split(',') : [];
        
        // Render existing tags
        tags.forEach(tag => addTag(container, tag.trim()));
        
        // Handle new tags
        input.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ',') {
                e.preventDefault();
                const value = input.value.trim();
                if (value) {
                    addTag(container, value);
                    input.value = '';
                    updateTagsInput(container);
                }
            }
        });
    }
    
    function addTag(container, value) {
        const tag = document.createElement('span');
        tag.className = 'tag';
        tag.innerHTML = `
            ${value}
            <span class="tag-remove" onclick="removeTag(this)">×</span>
        `;
        container.insertBefore(tag, container.querySelector('.tag-input'));
    }
    
    function removeTag(element) {
        const tag = element.parentElement;
        const container = tag.parentElement;
        tag.remove();
        updateTagsInput(container);
    }
    
    function updateTagsInput(container) {
        const tags = Array.from(container.querySelectorAll('.tag'))
            .map(tag => tag.textContent.replace('×', '').trim());
        const hiddenInput = container.parentElement.querySelector('input[type="hidden"]');
        hiddenInput.value = tags.join(',');
    }
    
    // Form validation and submission
    function collectFormValues() {
        const values = {};
        const form = document.getElementById('template-form');
        
        // Regular inputs
        form.querySelectorAll('input[name], select[name], textarea[name]').forEach(field => {
            if (field.type !== 'hidden' || field.name === 'etiketler') {
                values[field.name] = field.value;
            }
        });
        
        return values;
    }
    
    function validateForm() {
        const form = document.getElementById('template-form');
        let isValid = true;
        
        form.querySelectorAll('[required]').forEach(field => {
            const formGroup = field.closest('.form-group');
            if (!field.value.trim()) {
                formGroup.classList.add('error');
                isValid = false;
            } else {
                formGroup.classList.remove('error');
            }
        });
        
        return isValid;
    }
    
    function handlePreview() {
        if (!validateForm()) {
            showFormValidationErrors();
            return;
        }

        formValues = collectFormValues();

        // Show loading state on preview
        showLoadingState(elements.taskPreview, 'Önizleme hazırlanıyor...');
        showStep('preview');

        vscode.postMessage({ command: 'previewTask', values: formValues });
    }

    function handleCreate() {
        if (!validateForm()) {
            showFormValidationErrors();
            return;
        }

        formValues = collectFormValues();

        // Disable create button and show loading
        const createBtn = document.getElementById('btn-create') || document.getElementById('btn-confirm-create');
        if (createBtn) {
            createBtn.disabled = true;
            createBtn.innerHTML = '⏳ Oluşturuluyor...';
        }

        vscode.postMessage({ command: 'createTask', values: formValues });
    }

    function showFormValidationErrors() {
        // Highlight all invalid fields and scroll to first error
        const form = document.getElementById('template-form');
        const firstError = form.querySelector('.form-group.error input, .form-group.error textarea, .form-group.error select');

        if (firstError) {
            firstError.focus();
            firstError.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }

        // Show a notification
        showNotification('Lütfen tüm gerekli alanları doldurun', 'error');
    }

    function showNotification(message, type = 'info') {
        // Create notification element
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.innerHTML = `
            <div class="notification-content">
                <span class="notification-icon">${type === 'error' ? '❌' : type === 'success' ? '✅' : 'ℹ️'}</span>
                <span class="notification-message">${message}</span>
                <button class="notification-close" onclick="this.parentElement.parentElement.remove()">×</button>
            </div>
        `;

        // Add to page
        document.body.appendChild(notification);

        // Auto-remove after 5 seconds
        setTimeout(() => {
            if (notification.parentElement) {
                notification.remove();
            }
        }, 5000);
    }
    
    // Favorites handling
    function toggleFavorite(templateId, event) {
        event.stopPropagation();

        const favoriteBtn = event.target;
        const wasFavorite = isFavorite(templateId);

        if (wasFavorite) {
            removeFromFavorites(templateId);
            favoriteBtn.classList.remove('active');
            favoriteBtn.textContent = '☆';
        } else {
            addToFavorites(templateId);
            favoriteBtn.classList.add('active');
            favoriteBtn.textContent = '⭐';
        }
    }
    
    function isFavorite(templateId) {
        const favorites = getFavorites();
        return favorites.includes(templateId);
    }

    function getFavorites() {
        try {
            const stored = localStorage.getItem('gorev-template-favorites');
            return stored ? JSON.parse(stored) : [];
        } catch {
            return [];
        }
    }

    function saveFavorites(favorites) {
        try {
            localStorage.setItem('gorev-template-favorites', JSON.stringify(favorites));
        } catch (error) {
            console.error('Failed to save favorites:', error);
        }
    }

    function addToFavorites(templateId) {
        const favorites = getFavorites();
        if (!favorites.includes(templateId)) {
            favorites.push(templateId);
            saveFavorites(favorites);
            showNotification('Şablon favorilere eklendi', 'success');
            return true;
        }
        return false;
    }

    function removeFromFavorites(templateId) {
        const favorites = getFavorites();
        const index = favorites.indexOf(templateId);
        if (index > -1) {
            favorites.splice(index, 1);
            saveFavorites(favorites);
            showNotification('Şablon favorilerden çıkarıldı', 'info');
            return true;
        }
        return false;
    }
    
    // Message handling
    window.addEventListener('message', event => {
        const message = event.data;
        
        switch (message.command) {
            case 'templatesLoaded':
                // Check if we need to filter for favorites
                const activeTab = document.querySelector('.category-tab.active');
                if (activeTab && activeTab.dataset.category === 'favorites') {
                    const favoriteIds = getFavorites();
                    const favoriteTemplates = message.templates.filter(t => favoriteIds.includes(t.id));
                    renderTemplates(favoriteTemplates);
                } else {
                    renderTemplates(message.templates);
                }
                break;
            case 'searchResults':
                renderTemplates(message.templates);
                break;
            case 'favoritesLoaded':
                renderTemplates(message.templates);
                break;
                
            case 'templateSelected':
                renderForm(message.template);
                break;
                
            case 'previewGenerated':
                elements.taskPreview.innerHTML = marked.parse(message.preview);
                showStep('preview');
                break;
                
            case 'validationError':
                message.fields.forEach(fieldName => {
                    const field = document.querySelector(`[name="${fieldName}"]`);
                    if (field) {
                        field.closest('.form-group').classList.add('error');
                    }
                });
                break;
                
            case 'favoriteAdded':
                const btn = document.querySelector(
                    `.template-card[data-template-id="${message.templateId}"] .template-favorite`
                );
                if (btn) {
                    btn.classList.add('active');
                    btn.textContent = '⭐';
                }
                break;
        }
    });
    
    // Utility functions
    function debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }
    
    // Marked.js will be loaded from local bundle
    // The library is already included in the WebView bundle
    
    // Initialize when ready
    init();
})();

// Global functions for inline event handlers
function toggleFavorite(templateId, event) {
    event.stopPropagation();

    const favoriteBtn = event.target;
    const wasFavorite = favoriteBtn.classList.contains('active');

    if (wasFavorite) {
        const favorites = JSON.parse(localStorage.getItem('gorev-template-favorites') || '[]');
        const index = favorites.indexOf(templateId);
        if (index > -1) {
            favorites.splice(index, 1);
            localStorage.setItem('gorev-template-favorites', JSON.stringify(favorites));
        }
        favoriteBtn.classList.remove('active');
        favoriteBtn.textContent = '☆';
    } else {
        const favorites = JSON.parse(localStorage.getItem('gorev-template-favorites') || '[]');
        if (!favorites.includes(templateId)) {
            favorites.push(templateId);
            localStorage.setItem('gorev-template-favorites', JSON.stringify(favorites));
        }
        favoriteBtn.classList.add('active');
        favoriteBtn.textContent = '⭐';
    }
}

function removeTag(element) {
    const tag = element.parentElement;
    const container = tag.parentElement;
    tag.remove();
    
    // Update hidden input
    const tags = Array.from(container.querySelectorAll('.tag'))
        .map(tag => tag.textContent.replace('×', '').trim());
    const hiddenInput = container.parentElement.querySelector('input[type="hidden"]');
    hiddenInput.value = tags.join(',');
}