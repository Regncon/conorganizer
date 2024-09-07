'use client';
import type { EventDay, EventDays } from '$app/(public)/page';
import { useSetCustomCssVariable } from '$lib/hooks/useSetCustomCssVariable';
import { Box, Typography, Link, type SxProps, Divider } from '@mui/material';
import { Fragment, useState } from 'react';

type Props = {
    eventDays: EventDays;
    locationHash: EventDay;
};
const sxDayTypography: SxProps = {
    maxWidth: '5rem',
    minHeight: '4rem',
    display: 'grid',
    placeItems: 'center',
    paddingInline: '0.5em',
};

const DaysHeader = ({ eventDays, locationHash }: Props) => {
    const ref = useSetCustomCssVariable({ '--scroll-margin-top': 'height' });
    return (
        <>
            <Box
                component="header"
                sx={{
                    position: 'sticky',
                    top: 'var(--app-bar-height)',
                    backgroundColor: 'background.paper',
                    padding: '0.5rem',
                    zIndex: 1,
                }}
                ref={ref}
            >
                <Box sx={{ display: 'grid', placeContent: 'center', marginInline: '2rem', marginBlock: '0.5rem' }}>
                    FILTER
                </Box>
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(9, max-content)',
                        placeContent: 'center',
                        placeItems: 'center',
                    }}
                >
                    {Object.values(eventDays).map((day, i) => (
                        <Fragment key={day}>
                            {i === 0 && (
                                <Divider
                                    orientation="vertical"
                                    sx={{ borderColor: locationHash === day ? 'secondary.main' : '' }}
                                />
                            )}
                            <Box
                                sx={{
                                    backgroundColor: locationHash === day ? 'secondary.main' : '',
                                    transition: 'background-color 0.5s ease-in-out;',
                                }}
                            >
                                <Typography
                                    key={day}
                                    component={Link}
                                    href={`#${day}`}
                                    variant="h5"
                                    sx={sxDayTypography}
                                >
                                    {day}
                                </Typography>
                            </Box>
                            <Divider
                                orientation="vertical"
                                sx={{ borderColor: locationHash === day ? 'secondary.main' : '' }}
                            />
                        </Fragment>
                    ))}
                </Box>
            </Box>
        </>
    );
};

export default DaysHeader;
