import { Box, Button, Card, CardContent, Paper, Typography } from '@mui/material';
import CardBase from '$app/(authorized)/dashboard/CardBase';

export default function Ticket() {
    return (
        <Paper elevation={2} sx={{ width: '18rem', marginBottom: '2rem' }}>
            <Typography> Festivalpass+ </Typography>
            <Box sx={{ marginBlock: '0.75rem' }}>
                <Typography> John Doe </Typography>
                <Typography>johndoe@gmail.kil </Typography>
            </Box>
            <Typography> Billet nummer: 420P0GCH4MP </Typography>
            <Typography>Status: betalt </Typography>
            <Button variant="contained" color="primary">
                Tildel meg
            </Button>
            <Button variant="contained" color="primary">
                Tildel andre
            </Button>
        </Paper>
    );
}
