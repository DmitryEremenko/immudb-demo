import axios from 'axios';

export interface Account {
  account_number: string;
  account_name: string;
  iban: string;
  address: string;
  amount: string;
  type: string;
}

export const getAccounts = (): Promise<{ data: { accounts: Account[] } }> =>
  axios.get(`${import.meta.env.VITE_API_URL}/accounts`);

export const postAccount = (newAccount: Account) =>
  axios.post(`${import.meta.env.VITE_API_URL}/accounts`, newAccount);
