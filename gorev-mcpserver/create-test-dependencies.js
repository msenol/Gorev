#!/usr/bin/env node

// Simple script to create test tasks with dependencies for demonstration
const { spawn } = require('child_process');

async function createTask(title, description, priority) {
    return new Promise((resolve, reject) => {
        const params = {
            baslik: title,
            aciklama: description,
            oncelik: priority || 'orta'
        };
        
        console.log(`Creating task: ${title}`);
        // This would use the MCP client to create tasks
        // For now, just return a mock ID
        setTimeout(() => {
            resolve(`task-${Math.random().toString(36).substr(2, 9)}`);
        }, 100);
    });
}

async function createDependency(sourceId, targetId) {
    return new Promise((resolve, reject) => {
        console.log(`Creating dependency: ${sourceId} -> ${targetId}`);
        // This would use the MCP client to create dependency
        setTimeout(() => {
            resolve(true);
        }, 50);
    });
}

async function main() {
    console.log('Creating test tasks with dependencies...');
    
    // Create a chain of dependent tasks
    const task1Id = await createTask('Foundation Task', 'This is the foundation task that everything else depends on', 'yuksek');
    const task2Id = await createTask('Backend Setup', 'Setup the backend API - depends on Foundation', 'yuksek');
    const task3Id = await createTask('Frontend Development', 'Create the frontend UI - depends on Backend', 'orta');
    const task4Id = await createTask('Integration Testing', 'Test the integration - depends on Frontend', 'orta');
    const task5Id = await createTask('Deployment', 'Deploy to production - depends on Testing', 'dusuk');
    
    // Create dependency chain: task1 -> task2 -> task3 -> task4 -> task5
    await createDependency(task1Id, task2Id);
    await createDependency(task2Id, task3Id);
    await createDependency(task3Id, task4Id);
    await createDependency(task4Id, task5Id);
    
    // Create a parallel branch
    const task6Id = await createTask('Documentation', 'Write documentation - depends on Foundation', 'dusuk');
    const task7Id = await createTask('User Training', 'Train users - depends on Documentation and Frontend', 'dusuk');
    
    await createDependency(task1Id, task6Id);
    await createDependency(task6Id, task7Id);
    await createDependency(task3Id, task7Id);
    
    console.log('\n=== Test Dependencies Created ===');
    console.log('Dependency Chain:');
    console.log('1. Foundation Task → Backend Setup → Frontend Development → Integration Testing → Deployment');
    console.log('2. Foundation Task → Documentation → User Training');
    console.log('3. Frontend Development → User Training');
    console.log('\nWhat you should see in VS Code extension:');
    console.log('');
    console.log('TreeView badges:');
    console.log('- Foundation Task: [← 2] (2 tasks depend on this)');
    console.log('- Backend Setup: [🔗1] [← 1] (depends on 1, 1 depends on this)');
    console.log('- Frontend Development: [🔗1] [← 2] (depends on 1, 2 depend on this)');
    console.log('- Integration Testing: [🔗1] [← 1] (depends on 1, 1 depends on this)');
    console.log('- Deployment: [🔗1] (depends on 1)');
    console.log('- Documentation: [🔗1] [← 1] (depends on 1, 1 depends on this)');
    console.log('- User Training: [🔗2] (depends on 2)');
    console.log('');
    console.log('Task Detail Panel:');
    console.log('- Each task will show a "🔗 Bağımlılıklar" section');
    console.log('- "Bu görev için beklenen görevler" shows what this task depends on');
    console.log('- "Bu göreve bağımlı görevler" shows what depends on this task');
    console.log('- Dependency status indicators: ✅ completed, ⏳ pending');
}

main().catch(console.error);