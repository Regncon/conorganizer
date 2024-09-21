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
import Fuse from 'fuse.js'; // Import Fuse.js

type Props = {
    tickets: EventTicket[];
};

const TicketList = ({ tickets }: Props) => {
    const dinnerTicketId = 157059;

    // State to hold the search query
    const [searchQuery, setSearchQuery] = useState<string>('');

    // State to hold the filtered tickets
    const [filteredTickets, setFilteredTickets] = useState<EventTicket[]>(tickets);

    // Fuse.js options
    const fuseOptions = {
        keys: ['crm.first_name', 'crm.last_name', 'crm.email', 'order_id'], // Fields to search in
        threshold: 0.3, // Adjust to control fuzziness (lower is stricter)
    };

    useEffect(() => {
        const fuse = new Fuse(
            tickets.filter((ticket) => ticket.category_id !== dinnerTicketId), // Filter out dinner tickets
            fuseOptions
        );

        // Perform fuzzy search if searchQuery is not empty, else show all
        const result = searchQuery ? fuse.search(searchQuery).map((res) => res.item) : tickets;
        setFilteredTickets(result);
    }, [searchQuery, tickets]); // Re-run when search query or tickets change

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
