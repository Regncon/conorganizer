'use client';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import { Button, Container, InputAdornment, Paper, TextField } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import {
	signInAndCreateCookie,
	signOutAndDeleteCookie,
	singUpAndCreateCookie,
	type RegisterDetails,
} from '$lib/firebase/firebase';
import PasswordTextField from '../login/PasswordTextField';
import { useEffect, useRef } from 'react';

const Register = () => {
	const passwordRef = useRef<HTMLInputElement>(null);
	const confirmPasswordRef = useRef<HTMLInputElement>(null);
	const formRef = useRef<HTMLFormElement>(null);
	const emailRegExp = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+.[a-zA-Z]{2,4}$/;

	return (
		<Container component={Paper} fixed maxWidth="xl" sx={{ height: '70dvh' }}>
			<Grid2
				ref={formRef}
				component="form"
				container
				sx={{ placeContent: 'center', height: '100%', flexDirection: 'column', gap: '1rem' }}
				onSubmit={(e) => {
					e.preventDefault();
					const { password, confirm, email } = Object.fromEntries(
						new FormData(e.target as HTMLFormElement)
					) as RegisterDetails;

					if (password !== confirm) {
						console.log('no match', password !== confirm);
					}

					if (emailRegExp.test(email)) {
						singUpAndCreateCookie(e);
					}
				}}
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
						pattern: emailRegExp.source,
						title: 'epost@example.com',
					}}
				/>
				<PasswordTextField autoComplete="new-password" ref={passwordRef} />
				<PasswordTextField
					autoComplete="new-password"
					label="bekreft passord"
					name="confirm"
					ref={confirmPasswordRef}
				/>
				<Button type="submit">Log inn</Button>
			</Grid2>
		</Container>
	);
};

export default Register;
