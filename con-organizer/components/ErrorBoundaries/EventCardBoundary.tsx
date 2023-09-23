import { Alert, Button, Card, Typography } from '@mui/material';
type Props = {
    error: Error;
    resetErrorBoundary: () => void;
};
const EventCardBoundary = ({ error, resetErrorBoundary }: Props) => {
    return (
        <Card
            sx={{
                backgroundColor: 'primary.light',
                width: { xs: '100vw', md: '500px' },
                minHeight: '10em',
                display: 'grid',
                p: '1em',
                gridTemplateRows: '1fr auto 1fr auto 1fr',
            }}
        >
            <div></div>
            <Alert severity="error" sx={{ maxHeight: '11em', overflow: 'scroll' }}>
                <Typography variant="h6">Det har skjedd en feil:</Typography>
                <Typography variant="body1">{error.message}</Typography>
            </Alert>
            <div></div>
            <Button variant="contained" sx={{ backgroundColor: 'primary.main' }} onClick={resetErrorBoundary}>
                Pr&oslash;v igjen
            </Button>
            <div></div>
        </Card>
    );
};

export default EventCardBoundary;
