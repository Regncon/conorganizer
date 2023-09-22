import { Box, Button, Typography } from '@mui/material';
type Props = {
    error: Error;
    resetErrorBoundary: () => void;
};
const EventCardBoundary = ({ error, resetErrorBoundary }: Props) => {
    return (
        <Box sx={{ margin: '0 auto', backgroundColor: 'primary.light' }}>
            <Typography variant="h6">Det har skjedd en feil: {error.message}</Typography>
            <Button variant="contained" sx={{ backgroundColor: 'primary.main' }} onClick={resetErrorBoundary}>
                Pr&oslash;v igjen
            </Button>
        </Box>
    );
};

export default EventCardBoundary;
