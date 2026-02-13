export type ApprovalStatus = 'PENDING' | 'APPROVED' | 'REJECTED';

export interface IT03Doc {
  id: number;
  title: string;      // เอกสาร
  owner: string;      // ผู้ขอ (mock)
  status: ApprovalStatus;

  approvedReason?: string;
  rejectedReason?: string;
  updatedAt: string;
}
