'use client';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import { Button, Container, InputAdornment, Paper, TextField } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import PasswordTextField from './PasswordTextField';
import { signInAndCreateCookie, signOutAndDeleteCookie } from '$lib/firebase/firebase';
import { setSessionCookie } from './action';

const Login = () => {
	return (
		<Container component={Paper} fixed maxWidth="xl" sx={{ height: '70dvh' }}>
			<Grid2
				component="form"
				container
				sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
				onSubmit={signInAndCreateCookie}
			>
				<TextField
					type="email"
					name="email"
					autoComplete="email"
					label="e-post"
					variant="outlined"
					required
					InputProps={{
						endAdornment: (
							<InputAdornment position="end">
								<AccountCircleIcon />
							</InputAdornment>
						),
					}}
					inputProps={{
						pattern: '[a-z0-9._%+-]+@[a-z0-9.-]+.[a-z]{2,4}$',
					}}
				/>
				<PasswordTextField />
				<Button type="submit">Log inn</Button>
			</Grid2>
			<Button onClick={signOutAndDeleteCookie}>logg ut</Button>
		</Container>
	);
};

export default Login;
