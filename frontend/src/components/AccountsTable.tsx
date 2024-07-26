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
          {data?.data?.revisions?.map(({ document }) => (
            <TableRow key={document.account_number}>
              <TableCell>{document.account_number}</TableCell>
              <TableCell>{document.account_name}</TableCell>
              <TableCell>{document.iban}</TableCell>
              <TableCell>{document.address}</TableCell>
              <TableCell>{document.amount}</TableCell>
              <TableCell>{document.type}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
