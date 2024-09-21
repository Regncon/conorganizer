'use client';
import {
    Box,
    Card,
    CardActions,
    CardContent,
    FormControl,
    Input,
    InputAdornment,
    InputLabel,
    Typography,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import React, { useEffect, useState } from 'react';
import { EventTicket } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';

type Props = {
    tickets: EventTicket[];
};

const TicketList = ({ tickets }: Props) => {
    // Fetch tickets from the action
    const dinnerTicketId = 157059;

    // State to hold the search query
    const [searchQuery, setSearchQuery] = useState<string>('');

    // State to hold the filtered tickets
    const [filteredTickets, setFilteredTickets] = useState<EventTicket[]>([]);

    useEffect(() => {
        // Filter tickets based on the search query
        const filtered = tickets
            ?.filter(
                (ticket) =>
                    ticket.category_id !== dinnerTicketId &&
                    (ticket.crm.first_name.toLowerCase() + ' ' + ticket.crm.last_name.toLowerCase()).includes(
                        searchQuery.toLowerCase()
                    )
            )
            .sort((a, b) => a.crm.first_name.localeCompare(b.crm.first_name));

        setFilteredTickets(filtered || []);
    }, [searchQuery, tickets]); // Re-run the filter whenever searchQuery or tickets change

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
                            <InputLabel htmlFor="input-with-icon-adornment">Search for a ticket</InputLabel>
                            <Input
                                id="input-with-icon-adornment"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
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
                            <Typography>Ticket: {ticket.order_id}</Typography>
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

export default TicketList;
