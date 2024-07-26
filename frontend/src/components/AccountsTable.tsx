import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { getAccounts } from '../api';

export const AccountsTable = () => {
  const { data } = useQuery({ queryKey: ['accounts'], queryFn: getAccounts });

  return (
    <TableContainer component={Paper} sx={{ mt: '1rem' }}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Account Number</TableCell>
            <TableCell>Account Name</TableCell>
            <TableCell>IBAN</TableCell>
            <TableCell>Address</TableCell>
            <TableCell>Amount</TableCell>
            <TableCell>Type</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data?.data?.accounts?.map(
            (account: {
              account_number: string;
              account_name: string;
              iban: string;
              address: string;
              amount: string;
              type: string;
            }) => (
              <TableRow key={account.account_number}>
                <TableCell>{account.account_number}</TableCell>
                <TableCell>{account.account_name}</TableCell>
                <TableCell>{account.iban}</TableCell>
                <TableCell>{account.address}</TableCell>
                <TableCell>{account.amount}</TableCell>
                <TableCell>{account.type}</TableCell>
              </TableRow>
            ),
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
