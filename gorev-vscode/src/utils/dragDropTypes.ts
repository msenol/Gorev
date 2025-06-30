import { Gorev } from '../models/gorev';
import { GorevDurum, GorevOncelik } from '../models/common';

/**
 * Drag & Drop için veri transfer türleri
 */
export enum DragDataType {
    Task = 'application/vnd.gorev.task',
    Tasks = 'application/vnd.gorev.tasks',
    Status = 'application/vnd.gorev.status',
    Priority = 'application/vnd.gorev.priority',
    Project = 'application/vnd.gorev.project'
}

/**
 * Drag edilen görev verisi
 */
export interface TaskDragData {
    type: DragDataType.Task;
    task: Gorev;
    sourceGroupKey?: string;
}

/**
 * Çoklu görev drag verisi
 */
export interface TasksDragData {
    type: DragDataType.Tasks;
    tasks: Gorev[];
    sourceGroupKey?: string;
}

/**
 * Drop hedefi türleri
 */
export enum DropTargetType {
    StatusGroup = 'status-group',
    PriorityGroup = 'priority-group',
    ProjectGroup = 'project-group',
    Task = 'task',
    EmptyArea = 'empty-area'
}

/**
 * Drop hedefi bilgileri
 */
export interface DropTarget {
    type: DropTargetType;
    groupKey?: string;
    targetTask?: Gorev;
    newStatus?: GorevDurum;
    newPriority?: GorevOncelik;
    newProjectId?: string;
    position?: 'before' | 'after';
}

/**
 * Drag & Drop işlem sonucu
 */
export interface DragDropResult {
    success: boolean;
    updatedTasks?: Gorev[];
    error?: string;
}

/**
 * Drag & Drop konfigürasyonu
 */
export interface DragDropConfig {
    allowTaskMove: boolean;
    allowStatusChange: boolean;
    allowPriorityChange: boolean;
    allowProjectMove: boolean;
    allowDependencyCreate: boolean;
    allowParentChange: boolean;
    showDropIndicator: boolean;
    animateOnDrop: boolean;
}

/**
 * Drop zone görselleştirme
 */
export interface DropZoneVisual {
    showBefore: boolean;
    showAfter: boolean;
    showInside: boolean;
    highlightColor: string;
    borderStyle: 'solid' | 'dashed' | 'dotted';
}

/**
 * Drag feedback görselleştirme
 */
export interface DragFeedback {
    opacity: number;
    cursorStyle: string;
    showGhost: boolean;
    showBadge: boolean;
    badgeText?: string;
}