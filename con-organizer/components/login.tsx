'use client';

import { MouseEvent, useState } from 'react';
import Visibility from '@mui/icons-material/Visibility';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import { Alert } from '@mui/material';
import { Card, CardMedia } from '@mui/material';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';
import { signInWithEmailAndPassword } from 'firebase/auth';
import { auth } from '../lib/firebase';

const Login = (props: any) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [success, setSuccess] = useState('');
    const [error, setError] = useState('');
    const { setChoice } = props;
    const [showPassword, setShowPassword] = useState(false);
    const handleClickShowPassword = () => setShowPassword((show) => !show);
    const handleMouseDownPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };
    const login = (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        signInWithEmailAndPassword(auth, email, password)
            .then(() => {
                resetInput();
                setSuccess('Logget inn! Ett øyeblikk...');
                setError('');
                setChoice('');
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
        <Card>
            <Box p={4} maxWidth={400} display="grid" justifyItems="center" gap={2}>
                <CardMedia
                    component="img"
                    image="/img/regnconlogony.png"
                    title="årets regncondrage"
                    sx={{ p: '1rem' }}
                />
                <form action={''}>
                    <TextField
                        label="e-post"
                        value={email}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                        fullWidth
                        sx={{ margin: '.3rem 0' }}
                    />
                    <TextField
                        sx={{ margin: '.3rem 0' }}
                        label="passord"
                        name="password"
                        value={password}
                        type={showPassword ? 'text' : 'password'}
                        InputProps={{
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton
                                        aria-label="veksle mellom synlig og skjult passord"
                                        onClick={handleClickShowPassword}
                                        onMouseDown={handleMouseDownPassword}
                                        edge="end"
                                    >
                                        {showPassword ? <VisibilityOff /> : <Visibility />}
                                    </IconButton>
                                </InputAdornment>
                            ),
                        }}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                        fullWidth
                    />
                    <Button
                        variant="contained"
                        color="primary"
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
        </Card>
    );
};

export default Login;
