import { Component, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { It03Service } from '../../services/it03.service';
import { IT03Doc } from '../../models/it03.model';

type TabKey = 'ALL' | 'PENDING' | 'APPROVED' | 'REJECTED';
type ModalMode = 'APPROVE' | 'REJECT' | null;

@Component({
  selector: 'app-it03',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './it03.html',
  styleUrls: ['./it03.scss']
})
export class It03Component {

  // ===== STATE =====
  docs = signal<IT03Doc[]>([]);
  activeTab = signal<TabKey>('PENDING');
  selectedId = signal<number | null>(null);
  modalMode = signal<ModalMode>(null);
  reason = signal('');

  constructor(private it03: It03Service) {
    this.loadDocs();
  }

  // ===== COMPUTED =====
  filteredDocs = computed(() => {
    const tab = this.activeTab();
    const all = this.docs();

    if (tab === 'ALL') return all;
    return all.filter(d => d.status === tab);
  });

  selectedDoc = computed(() =>
    this.docs().find(d => d.id === this.selectedId()) ?? null
  );

  canApproveReject = computed(() =>
    this.selectedDoc()?.status === 'PENDING'
  );

  // ===== LOAD DATA =====
  private mapTabToStatus(tab: TabKey): string | undefined {
    if (tab === 'ALL') return undefined;
    return tab;
  }

  loadDocs() {
    const status = this.mapTabToStatus(this.activeTab());

    this.it03.getDocs(status).subscribe({
      next: (res: any) => {
        const list = res ?? [];
        console.log(list);
        
        this.docs.set(list);
      },
      error: (err) => {
        console.error(err);
        this.docs.set([]);
      }
    });
  }


  // ===== UI ACTION =====
  setTab(tab: TabKey) {
    this.activeTab.set(tab);
    this.selectedId.set(null);
    this.closeModal();
    this.loadDocs();
  }

  toggleSelect(id: number) {
    this.selectedId.set(this.selectedId() === id ? null : id);
  }

  openApprove() {
    if (!this.canApproveReject()) return;
    this.modalMode.set('APPROVE');
    this.reason.set('');
  }

  openReject() {
    if (!this.canApproveReject()) return;
    this.modalMode.set('REJECT');
    this.reason.set('');
  }

  closeModal() {
    this.modalMode.set(null);
    this.reason.set('');
  }

  // ===== SUBMIT =====
  submit() {
    const doc = this.selectedDoc();
    const mode = this.modalMode();
    const reason = this.reason().trim();

    if (!doc || !mode) return;

    if (!reason) {
      alert('กรุณากรอกเหตุผล');
      return;
    }

    const req$ =
      mode === 'APPROVE'
        ? this.it03.approve(doc.id, reason)
        : this.it03.reject(doc.id, reason);

    req$.subscribe({
      next: () => {
        this.closeModal();
        this.selectedId.set(null);
        this.loadDocs();
      },
      error: (err) => {
        alert(err?.error?.message || 'ดำเนินการไม่สำเร็จ');
        this.closeModal();
        this.loadDocs();
      }
    });
  }

  // ===== UI HELPER =====
  statusText(s: string) {
    if (s === 'PENDING') return 'รออนุมัติ';
    if (s === 'APPROVED') return 'อนุมัติ';
    if (s === 'REJECTED') return 'ไม่อนุมัติ';
    return s;
  }
}
