import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { IT03Doc } from '../models/it03.model';

@Injectable({ providedIn: 'root' })
export class It03Service {
  private readonly _docs$ = new BehaviorSubject<IT03Doc[]>([
    { id: 1, title: 'เอกสาร 1', owner: 'aaaa', status: 'PENDING', updatedAt: new Date().toISOString() },
    { id: 2, title: 'เอกสาร 2', owner: 'bbbb', status: 'PENDING', updatedAt: new Date().toISOString() },
    { id: 3, title: 'เอกสาร 3', owner: 'cccc', status: 'APPROVED', updatedAt: new Date().toISOString(), approvedReason: 'ครบถ้วน' },
    { id: 4, title: 'เอกสาร 4', owner: 'dddd', status: 'REJECTED', updatedAt: new Date().toISOString(), rejectedReason: 'ข้อมูลไม่ครบ' },
  ]);

  docs$ = this._docs$.asObservable();

  approve(id: number, reason: string) {
    const next: IT03Doc[] = this._docs$.value.map(d =>
      d.id === id
        ? { ...d, status: 'APPROVED', approvedReason: reason, rejectedReason: undefined, updatedAt: new Date().toISOString() }
        : d
    );
    this._docs$.next(next);
  }

  reject(id: number, reason: string) {
    const next: IT03Doc[] = this._docs$.value.map(d =>
      d.id === id
        ? { ...d, status: 'REJECTED', rejectedReason: reason, approvedReason: undefined, updatedAt: new Date().toISOString() }
        : d
    );
    this._docs$.next(next);
  }

}
