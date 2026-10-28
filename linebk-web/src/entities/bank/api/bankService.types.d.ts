export interface UserInfoResponse {
  name: string;
  greetingMessage: string;
}

export interface AccountResponse {
  type: string;
  amount: string;
  currency: string;
  accountNumber: string;
  issuer: string;
  color: string;
  isMainAccount: boolean;
  progress?: number;
  flags?: string[];
}

export interface RecentTransactionResponse {
  name: string;
  image: string;
  isBank: boolean;
}

export interface DebitCardResponse {
  name: string;
  status: string;
  issuer: string;
  color: string;
  borderColor?: string;
  number?: string;
}

export interface BannerResponse {
  title: string;
  description: string;
  image: string;
}
