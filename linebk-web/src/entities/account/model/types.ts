export interface AccountModel {
  type: string;
  amount: number;
  currency: string;
  accountNumber: string;
  issuer: string;
  color: string;
  isMainAccount: boolean;
  progress?: number;
  flags?: string[];
}
