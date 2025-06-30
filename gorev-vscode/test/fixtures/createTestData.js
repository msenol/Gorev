// Test verisi oluşturma scripti
// Bu scripti VS Code Extension Development Host'ta Debug Console'da çalıştırabilirsiniz

async function createTestHierarchy() {
    const mcpClient = vscode.extensions.getExtension('msenol.gorev-vscode').exports.mcpClient;
    
    if (!mcpClient || !mcpClient.isConnected()) {
        console.error('MCP Client not connected!');
        return;
    }
    
    try {
        // Ana proje oluştur
        console.log('Creating test project...');
        const projectResult = await mcpClient.callTool('proje_olustur', {
            isim: 'Test Hiyerarşi Projesi',
            tanim: 'Subtask özelliklerini test etmek için'
        });
        
        // Aktif proje yap
        await mcpClient.callTool('proje_aktif_yap', {
            proje_id: projectResult.content[0].text.match(/ID: ([a-f0-9-]+)/)[1]
        });
        
        // Root görev 1
        console.log('Creating root task 1...');
        const root1 = await mcpClient.callTool('gorev_olustur', {
            baslik: 'Website Yenileme Projesi',
            aciklama: 'Şirket websitesinin komple yenilenmesi',
            oncelik: 'yuksek',
            etiketler: 'frontend,urgent'
        });
        const root1Id = root1.content[0].text.match(/ID: ([a-f0-9-]+)/)[1];
        
        // Root 1'in alt görevleri
        console.log('Creating subtasks for root 1...');
        const design = await mcpClient.callTool('gorev_alt_gorev_olustur', {
            parent_id: root1Id,
            baslik: 'UI/UX Tasarım',
            aciklama: 'Yeni tasarımın hazırlanması',
            oncelik: 'yuksek'
        });
        const designId = design.content[0].text.match(/ID: ([a-f0-9-]+)/)[1];
        
        const frontend = await mcpClient.callTool('gorev_alt_gorev_olustur', {
            parent_id: root1Id,
            baslik: 'Frontend Geliştirme',
            aciklama: 'React ile frontend implementasyonu',
            oncelik: 'orta'
        });
        const frontendId = frontend.content[0].text.match(/ID: ([a-f0-9-]+)/)[1];
        
        const backend = await mcpClient.callTool('gorev_alt_gorev_olustur', {
            parent_id: root1Id,
            baslik: 'Backend API',
            aciklama: 'REST API geliştirme',
            oncelik: 'orta'
        });
        
        // Design'ın alt görevleri (3. seviye)
        console.log('Creating sub-subtasks...');
        await mcpClient.callTool('gorev_alt_gorev_olustur', {
            parent_id: designId,
            baslik: 'Wireframe Hazırlama',
            aciklama: 'İlk taslaklar',
            oncelik: 'yuksek'
        });
        
        await mcpClient.callTool('gorev_alt_gorev_olustur', {
            parent_id: designId,
            baslik: 'Mockup Tasarım',
            aciklama: 'Detaylı tasarımlar',
            oncelik: 'orta'
        });
        
        // Bazı görevleri tamamla
        console.log('Completing some tasks...');
        await mcpClient.callTool('gorev_guncelle', {
            id: designId,
            durum: 'tamamlandi'
        });
        
        // Root görev 2
        console.log('Creating root task 2...');
        const root2 = await mcpClient.callTool('gorev_olustur', {
            baslik: 'Mobil Uygulama',
            aciklama: 'iOS ve Android uygulaması',
            oncelik: 'orta',
            etiketler: 'mobile'
        });
        const root2Id = root2.content[0].text.match(/ID: ([a-f0-9-]+)/)[1];
        
        // Root 2'nin alt görevi
        await mcpClient.callTool('gorev_alt_gorev_olustur', {
            parent_id: root2Id,
            baslik: 'iOS Geliştirme',
            aciklama: 'Swift ile iOS uygulaması',
            oncelik: 'orta'
        });
        
        // Bağımsız görev (drag & drop testi için)
        console.log('Creating standalone task...');
        await mcpClient.callTool('gorev_olustur', {
            baslik: 'Dokümantasyon Güncelleme',
            aciklama: 'Drag & drop ile taşınabilir',
            oncelik: 'dusuk',
            etiketler: 'docs'
        });
        
        console.log('Test hierarchy created successfully!');
        console.log('Refresh the tree view to see the hierarchy.');
        
        // Tree'yi yenile
        vscode.commands.executeCommand('gorev.refreshTasks');
        
    } catch (error) {
        console.error('Error creating test data:', error);
    }
}

// Fonksiyonu çalıştır
createTestHierarchy();