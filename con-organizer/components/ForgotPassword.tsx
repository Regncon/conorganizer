"use client";

import { useState } from 'react';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import { Alert, IconButton } from '@mui/material';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import { sendPasswordResetEmail } from 'firebase/auth';
import { auth } from '../lib/firebase';

const ForgotPassword = (props: any) => {
    const [email, setEmail] = useState('');
    const [success, setSuccess] = useState('');
    const [error, setError] = useState('');
    const { setChoice } = props;
    const [showPassword, setShowPassword] = useState(false);
    const handleClickShowPassword = () => setShowPassword((show) => !show);
    const handleMouseDownPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };

    const resetPwd = () => {
        sendPasswordResetEmail(auth, email)
            .then(() => {
                resetInput();
                setSuccess('Suksess! Vi har sendt deg en lenke for å skrive inn et nytt passord. Sjekk eposten din!');
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
            setError('Kunne ikke sende epost, fordi det skjedde en feil. Tekniske detaljer: ' + err.message);
        }
        setError('Klarte ikke sende epost, ta kontakt hvis problemet vedvarer! Tekniske detaljer: ' + err.message);
            });
    };

    const resetInput = () => {
        setEmail('');
    };
    return (
        <Box p={5} maxWidth={600} display={'grid'} justifyItems={'center'} gap={1}>
            <h1>Glemt/endre passord</h1>
            <TextField
                label="e-post"
                id="outlined-adornment-password"
                name="new-password"
                value={email}
                type={showPassword ? 'text' : 'password'}
                endAdornment={
                    <InputAdornment position="end">
                        <IconButton
                            aria-label="toggle password visibility"
                            onClick={handleClickShowPassword}
                            onMouseDown={handleMouseDownPassword}
                            edge="end"
                        >
                            {showPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                    </InputAdornment>
                }
                onChange={(e) => setEmail(e.target.value)}
                fullWidth
            />
            <Button variant="contained" size="large" fullWidth onClick={resetPwd}>
                Send!
            </Button>
            <Button variant="outlined" size="large" fullWidth onClick={() => setChoice('')}>
                Avbryt
            </Button>
            {!!success && <Alert severity="success">{success}</Alert>}
            {!!error && <Alert severity="error">{error}</Alert>}
        </Box>
    );
};

export default ForgotPassword;
