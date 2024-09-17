import React, { useState, useEffect } from 'react';
import { Accordion, AccordionDetails, AccordionSummary, Box, IconButton, Paper, Typography } from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ArrowDropUpIcon from '@mui/icons-material/ArrowDropUp';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { ConEvent } from '$lib/types';

type Props = {
    id: string;
    pulje: string;
    allEvents: ConEvent[];
    disabled: boolean;
};

const Ordering = ({ id, pulje, allEvents, disabled }: Props) => {
    const [expanded, setExpanded] = useState<boolean>(false);
    useEffect(() => {
        if (disabled) {
            setExpanded(false);
        }
    }, [disabled]);

    const handleAccordionChange = () => {
        if (!disabled) {
            setExpanded((prevExpanded) => !prevExpanded);
        }
    };

    return (
        <Paper elevation={3}>
            <Accordion expanded={expanded} onChange={handleAccordionChange} disabled={disabled}>
                <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel1-content" id="panel1-header">
                    Sortering rekkefølge
                </AccordionSummary>
                <AccordionDetails>
                    <Typography variant={'h3'} sx={{ textAlign: 'center' }}>
                        Pulje: {pulje}
                    </Typography>
                    {!allEvents ?
                        <Typography variant={'h3'} sx={{ textAlign: 'center' }}>
                            Laster inn eventer...
                        </Typography>
                        : <>
                            {allEvents.map((event: ConEvent) => (
                                <Paper
                                    key={event.id}
                                    elevation={4}
                                    sx={{
                                        padding: '1rem',
                                        marginBottom: '1rem',
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        backgroundColor: event.id === id ? 'primary.light' : '',
                                    }}
                                >
                                    <Typography component={'span'}>{event.title}</Typography>
                                    <Box
                                        sx={{
                                            display: 'inline-block',
                                        }}
                                    >
                                        <IconButton>
                                            <ArrowDropUpIcon />
                                        </IconButton>
                                        <IconButton>
                                            <ArrowDropDownIcon />
                                        </IconButton>
                                    </Box>
                                </Paper>
                            ))}
                        </>
                    }
                </AccordionDetails>
            </Accordion>
        </Paper>
    );
};
export default Ordering;
