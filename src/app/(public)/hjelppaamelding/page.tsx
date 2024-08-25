import BackButton from '$app/(authorized)/BackButton';
import { Paper } from '@mui/material';
import { Metadata } from 'next';

export const metadata: Metadata = {
    title: 'Hjelp påmelding',
    description: 'Forklaring på hvordan pujepåmeldingen fungerer',
};

const HjelpPaameldingPage = async () => {
    return (
        <Paper>
            <BackButton />
            <h1>Hjelp påmelding</h1>
        </Paper>
    );
};
export default HjelpPaameldingPage;
