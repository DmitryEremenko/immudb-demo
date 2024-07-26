import axios from 'axios';

export interface Account {
  account_number: string;
  account_name: string;
  iban: string;
  address: string;
  amount: string;
  type: string;
}

export const getAccounts = (): Promise<{
  data: { revisions: { document: Account }[] };
}> => axios.get(`${import.meta.env.VITE_API_URL}/documents?page=1&perPage=100`);

export const postAccount = (newAccount: Account) =>
  axios.put(`${import.meta.env.VITE_API_URL}/document`, newAccount);
