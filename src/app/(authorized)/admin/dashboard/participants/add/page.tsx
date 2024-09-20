import {
    Box,
    Card,
    CardActions,
    CardContent,
    CardHeader,
    FormControl,
    Input,
    InputAdornment,
    InputLabel,
    Stack,
    Typography,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import { EventTicket, GetTicketsFromCheckIn } from '$app/(authorized)/my-profile/my-tickets/actions';

const AddParticipant = async () => {
    const tickets: EventTicket[] | undefined = await GetTicketsFromCheckIn();
    const dinnerTicketId = 157059;
    const filteredTickets = tickets?.filter((ticket) => ticket.category_id !== dinnerTicketId);
    return (
        <Box>
            <h1>Add Participant</h1>
            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 1fr))',
                    gap: '1rem',
                }}
            >
                <Card>
                    <CardActions>
                        <FormControl variant="standard">
                            <InputLabel htmlFor="input-with-icon-adornment">Søk etter bilett</InputLabel>
                            <Input
                                id="input-with-icon-adornment"
                                endAdornment={
                                    <InputAdornment position="start">
                                        <SearchIcon />
                                    </InputAdornment>
                                }
                            />
                        </FormControl>
                    </CardActions>
                </Card>
                {filteredTickets?.map((ticket) => (
                    <Card key={ticket.id}>
                        <CardContent>
                            <Typography>Bilett: {ticket.order_id}</Typography>
                            <Typography>{ticket.category}</Typography>
                            <Typography>
                                {ticket.crm.first_name} {ticket.crm.last_name}
                            </Typography>
                            <Typography>{ticket.crm.email}</Typography>
                        </CardContent>
                    </Card>
                ))}
            </Box>
        </Box>
    );
};
export default AddParticipant;
