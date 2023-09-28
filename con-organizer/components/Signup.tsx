'use client';

import { MouseEvent, useState } from 'react';
import { Visibility, VisibilityOff } from '@mui/icons-material';
import { Alert, Box, Button, Card, IconButton, InputAdornment, Link, TextField } from '@mui/material';
import { FirebaseError } from 'firebase/app';
import { sendEmailVerification, signInWithEmailAndPassword } from 'firebase/auth';
import { auth } from '@/lib/firebase';
import { GetLoginInfoResponse } from '@/models/enums';

type Props = {
    setChoice: (choice: string) => void;
};

const Signup = ({ setChoice }: Props) => {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [passwordConfirmation, setPasswordConfirmation] = useState<string>('');
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const [success, setSuccess] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [showAlert, setShowAlert] = useState<boolean>(false);

    const handleClickShowPassword = () => setShowPassword((show) => !show);
    const handleMouseDownPassword = (event: MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };
    const signUp = async (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        if (password !== passwordConfirmation) {
            setError('Passordene matsjer isje, prøv igjen');
            return;
        }
        const result = await fetch('/api/getlogininfo', { method: 'POST', body: JSON.stringify({ email, password }) });
        const res = await result.json();
        console.log(res);
        if (res.user === GetLoginInfoResponse.Created) {
            setShowAlert(false);
            try {
                const signedInUser = await signInWithEmailAndPassword(auth, email, password);
                await sendEmailVerification(signedInUser.user);
                setSuccess(
                    'Du skal ha fått en verifiserings e-post, sjekk søppelpost vist du ikkje ser den i innboksen'
                );
                setError('');
                console.log(signedInUser);
                auth.signOut();
            } catch (e) {
                const err = e as FirebaseError;
                console.error(e);
                if (err.code === 'auth/invalid-email') {
                    setError('Det ser ut som du skrev noe feil i epostadressen din, vennligst prøv igjen.');
                } else if (err.code === 'auth/too-many-requests') {
                    setError(
                        'Vi har midlertidig suspendert kontoen din på grunn av for mange påloggingsforsøk. Dette er for at hackere og botter ikke skal kunne gjette seg til passordet ditt. Vennligst vent litt før du prøver igjen.'
                    );
                } else if (err.code === 'auth/weak-password') {
                    setError('Passordet må være minst seks tegn langt');
                } else {
                    setError('Kunne ikke registrere deg fordi det skjedde en feil. Tekniske detaljer: ' + err.message);
                }
            }
        }
        if (res.user === GetLoginInfoResponse.Exists) {
            setError('Du har allerede laget en bruker gå til login');
            setShowAlert(false);
        }
        if (res.user === GetLoginInfoResponse.DontExist) {
            setError('Du må bruke samme epost som du brukte til og kjøpe billett/ eller du må kjøpe billet');
            setShowAlert(true);
        }
        return;
    };
    return (
        <Box>
            <Card>
                <Box p={4} maxWidth={400} display="grid" justifyItems="center" gap={2}>
                    <form action={''}>
                        <TextField
                            label="e-post"
                            value={email}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                            fullWidth
                            sx={{ margin: '.3rem 0' }}
                        />
                        <TextField
                            sx={{ margin: '-1px 0 .3rem 0' }}
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
                        <TextField
                            sx={{ margin: '-1px 0 .3rem 0' }}
                            label="passord (igjen)"
                            name="password"
                            value={passwordConfirmation}
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
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                setPasswordConfirmation(e.target.value)
                            }
                            fullWidth
                        />
                        <Button
                            variant="contained"
                            color="primary"
                            size="large"
                            type="submit"
                            fullWidth
                            onClick={(e) => signUp(e)}
                            sx={{ margin: '.3rem 0' }}
                        >
                            Lag konto
                        </Button>
                        <Button
                            variant="outlined"
                            size="large"
                            fullWidth
                            onClick={() => setChoice('')}
                            sx={{ margin: '-1px 0' }}
                        >
                            Avbryt
                        </Button>

                        {!!success && <Alert severity="success">{success}</Alert>}
                        {!!error && <Alert severity="error">{error}</Alert>}
                        {showAlert && (
                            <Alert severity="info">
                                OBS: Du kan ikke lage bruker uten &aring; ha kj&oslash;pt billett.&nbsp;
                                <Link href="https://www.regncon.no/kjop-billett-til-regncon-xxxi/">
                                    Kjøp billett her!
                                </Link>
                            </Alert>
                        )}
                    </form>
                </Box>
            </Card>
        </Box>
    );
};

export default Signup;
