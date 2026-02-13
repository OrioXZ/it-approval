import { Injectable } from '@angular/core';
import { IT03Doc } from '../models/it03.model';
import { HttpClient } from '@angular/common/http';
import { map, Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class It03Service {

  private baseUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) { }

  getDocs(status?: string) {
  const q = status ? `?status=${status}` : '';

  return this.http
    .get<any>(`${this.baseUrl}/requests${q}`)
    .pipe(
      map(res => (res.items ?? []).map((i: any) => ({
        id: i.id,
        title: i.title,
        owner: i.decided_by || '-',
        status: i.status_code,
        approvedReason: i.status_code === 'APPROVED' ? i.decided_reason : undefined,
        rejectedReason: i.status_code === 'REJECTED' ? i.decided_reason : undefined,
        updatedAt: i.updated_at
      })))
    );
}

  
  approve(id: number, reason: string) {
    return this.http.patch(`${this.baseUrl}/requests/${id}/status`, {
      statusCode: 'APPROVED',
      decidedReason: reason
    });
  }

  reject(id: number, reason: string) {
    return this.http.patch(`${this.baseUrl}/requests/${id}/status`, {
      statusCode: 'REJECTED',
      decidedReason: reason
    });
  }

}

