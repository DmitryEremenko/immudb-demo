import { Container, Typography, Box } from '@mui/material';
import { AccountsTable } from './components/AccountsTable';
import { AccountForm } from './components/AccountForm';

function App() {
  return (
    <Container>
      <Typography variant='h4' component='h1' gutterBottom>
        Accounting Information
      </Typography>
      <Box sx={{ display: 'grid', gridTemplateColumns: '0.5fr 1fr', gap: '2rem' }}>
        <AccountForm />
        <AccountsTable />
      </Box>
    </Container>
  );
}

export default App;
