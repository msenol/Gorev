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
            vscode.postMessage({ command: 'loadFavorites' });
        } else {
            loadTemplatesByCategory(category);
        }
    }
    
    function loadTemplatesByCategory(category) {
        vscode.postMessage({ 
            command: 'loadTemplates', 
            category: category || undefined 
        });
    }
    
    function renderTemplates(templateList) {
        templates = templateList;
        
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
    
    // Form handling
    function renderForm(template) {
        selectedTemplate = template;
        
        // Update header
        elements.templateName.textContent = template.isim;
        elements.templateDescription.textContent = template.tanim || '';
        
        // Build form fields
        const fieldsHtml = template.alanlar.map(field => {
            const fieldId = `field-${field.isim}`;
            const required = field.zorunlu ? 'required' : '';
            const value = formValues[field.isim] || field.varsayilan || '';
            
            switch (field.tur) {
                case 'text':
                    return `
                        <div class="form-group">
                            <label class="form-label ${required}" for="${fieldId}">${field.isim}</label>
                            <input type="text" id="${fieldId}" name="${field.isim}" 
                                   class="form-input" value="${value}" ${required}>
                            ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                            <div class="form-error">Bu alan zorunludur</div>
                        </div>
                    `;
                    
                case 'textarea':
                    return `
                        <div class="form-group">
                            <label class="form-label ${required}" for="${fieldId}">${field.isim}</label>
                            <textarea id="${fieldId}" name="${field.isim}" 
                                      class="form-textarea" ${required}>${value}</textarea>
                            ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                            <div class="form-error">Bu alan zorunludur</div>
                        </div>
                    `;
                    
                case 'select':
                    const options = field.secenekler || ['dusuk', 'orta', 'yuksek'];
                    return `
                        <div class="form-group">
                            <label class="form-label ${required}" for="${fieldId}">${field.isim}</label>
                            <select id="${fieldId}" name="${field.isim}" 
                                    class="form-select" ${required}>
                                ${options.map(opt => 
                                    `<option value="${opt}" ${value === opt ? 'selected' : ''}>${opt}</option>`
                                ).join('')}
                            </select>
                            ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                            <div class="form-error">Bu alan zorunludur</div>
                        </div>
                    `;
                    
                case 'date':
                    return `
                        <div class="form-group">
                            <label class="form-label ${required}" for="${fieldId}">${field.isim}</label>
                            <input type="date" id="${fieldId}" name="${field.isim}" 
                                   class="form-input" value="${value}" ${required}>
                            ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                            <div class="form-error">Bu alan zorunludur</div>
                        </div>
                    `;
                    
                case 'tags':
                    return `
                        <div class="form-group">
                            <label class="form-label ${required}" for="${fieldId}">${field.isim}</label>
                            <div class="tag-input-container" id="${fieldId}-container">
                                <input type="text" id="${fieldId}" class="tag-input" 
                                       placeholder="Etiket ekle ve Enter'a bas...">
                            </div>
                            <input type="hidden" name="${field.isim}" value="${value}">
                            ${field.aciklama ? `<div class="form-help">${field.aciklama}</div>` : ''}
                            <div class="form-error">Bu alan zorunludur</div>
                        </div>
                    `;
                    
                default:
                    return '';
            }
        }).join('');
        
        elements.formFields.innerHTML = fieldsHtml;
        
        // Initialize tag inputs
        document.querySelectorAll('.tag-input').forEach(input => {
            initializeTagInput(input);
        });
        
        // Show form step
        showStep('form-fields');
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
            return;
        }
        
        formValues = collectFormValues();
        vscode.postMessage({ command: 'previewTask', values: formValues });
    }
    
    function handleCreate() {
        if (!validateForm()) {
            return;
        }
        
        formValues = collectFormValues();
        vscode.postMessage({ command: 'createTask', values: formValues });
    }
    
    // Favorites handling
    function toggleFavorite(templateId, event) {
        event.stopPropagation();
        vscode.postMessage({ command: 'saveAsFavorite', templateId });
    }
    
    function isFavorite(templateId) {
        // This would be managed by the extension
        return false;
    }
    
    // Message handling
    window.addEventListener('message', event => {
        const message = event.data;
        
        switch (message.command) {
            case 'templatesLoaded':
            case 'searchResults':
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
    
    // Add marked.js for markdown parsing
    const script = document.createElement('script');
    script.src = 'https://cdn.jsdelivr.net/npm/marked/marked.min.js';
    document.head.appendChild(script);
    
    // Initialize when ready
    init();
})();

// Global functions for inline event handlers
function toggleFavorite(templateId, event) {
    event.stopPropagation();
    vscode.postMessage({ command: 'saveAsFavorite', templateId });
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