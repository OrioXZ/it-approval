import { Component, computed, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { It03Service } from '../../services/it03.service';
import { IT03Doc, ApprovalStatus } from '../../models/it03.model';

type TabKey = 'ALL' | 'PENDING' | 'APPROVED' | 'REJECTED';

@Component({
  selector: 'app-it03',
  imports: [CommonModule, FormsModule],
  templateUrl: './it03.html',
  styleUrl: './it03.scss',
})
export class It03 {
// data
  docs = signal<IT03Doc[]>([]);
  constructor(private it03: It03Service) {
    this.it03.docs$.subscribe(v => this.docs.set(v));
  }

  // tab
  activeTab = signal<TabKey>('ALL');

  filteredDocs = computed(() => {
    const tab = this.activeTab();
    const all = this.docs();
    if (tab === 'ALL') return all;
    return all.filter(d => d.status === tab);
  });

  // selection: เลือกได้ทีละรายการ (เหมือนรูปมี checkbox)
  selectedId = signal<number | null>(null);
  selectedDoc = computed(() => this.docs().find(d => d.id === this.selectedId()) ?? null);

  // modal state
  modalOpen = signal(false);
  modalMode = signal<'APPROVE' | 'REJECT' | null>(null);
  reason = signal('');

  // disable rule: ถ้าอนุมัติแล้ว "อนุมัติซ้ำไม่ได้" และรายการไม่ใช่รออนุมัติห้ามกดทั้งคู่
  canAct = computed(() => this.selectedDoc()?.status === 'PENDING');

  setTab(tab: TabKey) {
    this.activeTab.set(tab);
    this.selectedId.set(null);
    this.closeModal();
  }

  select(id: number) {
    this.selectedId.set(this.selectedId() === id ? null : id);
    this.closeModal();
  }

  openApprove() {
    if (!this.canAct()) return;
    this.modalMode.set('APPROVE');
    this.reason.set('');
    this.modalOpen.set(true);
  }

  openReject() {
    if (!this.canAct()) return;
    this.modalMode.set('REJECT');
    this.reason.set('');
    this.modalOpen.set(true);
  }

  submit() {
    const doc = this.selectedDoc();
    const mode = this.modalMode();
    const reason = this.reason().trim();

    if (!doc || !mode) return;
    if (!reason) return; // ให้ผู้ใช้กรอกเหตุผลตามโจทย์ :contentReference[oaicite:2]{index=2}

    if (mode === 'APPROVE') this.it03.approve(doc.id, reason);
    if (mode === 'REJECT') this.it03.reject(doc.id, reason);

    this.closeModal();
  }

  closeModal() {
    this.modalOpen.set(false);
    this.modalMode.set(null);
    this.reason.set('');
  }

  statusText(s: ApprovalStatus) {
    if (s === 'PENDING') return 'รออนุมัติ';
    if (s === 'APPROVED') return 'อนุมัติ';
    return 'ไม่อนุมัติ';
  }
}
