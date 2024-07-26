import { useState, useCallback } from 'react';
import {
  Container,
  Typography,
  TextField,
  Button,
  Select,
  MenuItem,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from '@mui/material';
import { Account, getAccounts, postAccount } from './api';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

function App() {
  const [newAccount, setNewAccount] = useState<Account>({
    account_number: '',
    account_name: '',
    iban: '',
    address: '',
    amount: '0',
    type: 'sending',
  });
  const queryClient = useQueryClient();

  const { data } = useQuery({ queryKey: ['accounts'], queryFn: getAccounts });

  const { mutate } = useMutation({
    mutationFn: postAccount,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] });
    },
  });

  const handleInputChange = useCallback(
    (e: { target: { name: string; value: string } }) => {
      setNewAccount({ ...newAccount, [e.target.name]: e.target.value });
    },
    [newAccount],
  );

  const handleSubmit = () => {
    mutate(newAccount);
  };

  return (
    <Container maxWidth='md'>
      <Typography variant='h4' component='h1' gutterBottom>
        Accounting Information
      </Typography>
      <form onSubmit={handleSubmit} style={{ marginBottom: '2rem' }}>
        <TextField
          fullWidth
          margin='normal'
          name='account_number'
          label='Account Number'
          value={newAccount.account_number}
          onChange={handleInputChange}
          required
        />
        <TextField
          fullWidth
          margin='normal'
          name='account_name'
          label='Account Name'
          value={newAccount.account_name}
          onChange={handleInputChange}
          required
        />
        <TextField
          fullWidth
          margin='normal'
          name='iban'
          label='IBAN'
          value={newAccount.iban}
          onChange={handleInputChange}
          required
        />
        <TextField
          fullWidth
          margin='normal'
          name='address'
          label='Address'
          value={newAccount.address}
          onChange={handleInputChange}
          required
        />
        <TextField
          fullWidth
          margin='normal'
          name='amount'
          label='Amount'
          type='number'
          value={newAccount.amount}
          onChange={handleInputChange}
          required
        />
        <Select
          fullWidth
          margin='none'
          name='type'
          value={newAccount.type}
          onChange={handleInputChange}
          required
        >
          <MenuItem value='sending'>Sending</MenuItem>
          <MenuItem value='receiving'>Receiving</MenuItem>
        </Select>
        <Button type='submit' variant='contained' color='primary' style={{ marginTop: '1rem' }}>
          Add Account
        </Button>
      </form>
      <Typography variant='h5' component='h2' gutterBottom>
        Accounts
      </Typography>
      <TableContainer component={Paper}>
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
    </Container>
  );
}

export default App;
