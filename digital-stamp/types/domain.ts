export type Customer = {
  id: string
  name: string
  phone: string
  created_at: string
  updated_at: string
}

export type Staff = {
  id: string
  name: string
  email: string
  created_at: string
  updated_at: string
}

export type StampCard = {
  id: string
  customer_id: string
  stamp_count: number
  required_stamps: number
  created_at: string
  updated_at: string
}

export type Reward = {
  id: string
  name: string
  description: string
  required_stamps: number
  active: boolean
  created_at: string
  updated_at: string
}

export type QRStatus = "ACTIVE" | "USED" | "EXPIRED" | "CANCELLED"
export type TransactionType = "STAMP_ADDED" | "STAMP_REVERSED" | "REWARD_REDEEMED"

export type StampQR = {
  id: string
  staff_id: string
  status: QRStatus
  expires_at: string
  used_at?: string
  used_by_customer_id?: string
  created_at: string
}

export type StampTransaction = {
  id: string
  customer_id: string
  staff_id?: string
  reward_id?: string
  stamp_qr_id?: string
  type: TransactionType
  stamp_delta: number
  created_at: string
}
