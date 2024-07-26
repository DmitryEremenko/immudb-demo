import { useState, useCallback } from 'react';
import { TextField, Button, Select, MenuItem } from '@mui/material';
import { Account, postAccount } from '../api';
import { useMutation, useQueryClient } from '@tanstack/react-query';

export const AccountForm = () => {
  const [newAccount, setNewAccount] = useState<Account>({
    account_number: '',
    account_name: '',
    iban: '',
    address: '',
    amount: '0',
    type: 'sending',
  });
  const queryClient = useQueryClient();

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
    <form onSubmit={handleSubmit}>
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
  );
};
