'use client';

import { MouseEvent, useState } from 'react';
import { Alert } from '@mui/material';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import { signInWithEmailAndPassword } from 'firebase/auth';
import { auth } from '../lib/firebase';

const Login = (props: any) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [success, setSuccess] = useState('');
    const [error, setError] = useState('');
    const { setChoice } = props;
    const login = (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        signInWithEmailAndPassword(auth, email, password)
            .then(() => {
                resetInput();
                setSuccess('Logget inn! Ett øyeblikk...');
                setError('');
            })
            .catch((err) => {
                console.error(err);
                if (err.code === 'auth/invalid-email') {
                    setError('Det ser ut som du skrev noe feil i epostadressen din, vennligst prøv igjen.');
                } else if (err.code === 'auth/too-many-requests') {
                    setError(
                        'Vi har midlertidig suspendert kontoen din på grunn av for mange påloggingsforsøk. Dette er for at hackere og botter ikke skal kunne gjette seg til passordet ditt. Vennligst vent litt før du prøver igjen.'
                    );
                } else if (err.code === 'auth/wrong-password') {
                    setError('Ser ut som du har skrevet noe feil i passordet ditt, vennligst prøv igjen.');
                } else {
                    setError('Kunne ikke logge deg inn, fordi det skjedde en feil. Tekniske detaljer: ' + err.message);
                }
            });
    };
    const resetInput = () => {
        setEmail('');
        setPassword('');
    };
    return (
        <Box p={5} maxWidth={600} display="grid" justifyItems="center" gap={2}>
            <img src="/img/regnconlogony.png" alt="årets regncondrage" width={200} />
            <form action={''}>
                <TextField
                    label="e-post"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    fullWidth
                    sx={{ margin: '.3rem 0' }}
                />
                <TextField
                    sx={{ margin: '.3rem 0' }}
                    label="passord"
                    name="password"
                    value={password}
                    type="password"
                    onChange={(e) => setPassword(e.target.value)}
                    fullWidth
                />
                <Button
                    variant="contained"
                    size="large"
                    type="submit"
                    fullWidth
                    onClick={login}
                    sx={{ margin: '.3rem 0' }}
                >
                    Logg inn
                </Button>
                <Button
                    variant="outlined"
                    size="large"
                    fullWidth
                    onClick={() => setChoice('')}
                    sx={{ margin: '.3rem 0' }}
                >
                    Avbryt
                </Button>
            </form>
            {!!success && <Alert severity="success">{success}</Alert>}
            {!!error && <Alert severity="error">{error}</Alert>}
        </Box>
    );
};

export default Login;
