'use client';
import { PoolName } from '$lib/enums';
import { useSetCustomCssVariable } from '$lib/hooks/useSetCustomCssVariable';
import { Box, Typography, Link, type SxProps, Divider } from '@mui/material';
import { Fragment } from 'react';
import { translatedDays } from '../lib/helpers/translation';

const sxDayTypography: SxProps = {
    maxWidth: '5.5rem',
    minHeight: '4rem',
    display: 'grid',
    placeItems: 'center',
    paddingInline: '0.5em',
    textAlign: 'center',
    transition: 'background-color 0.5s ease-in-out;',
};

type Props = {};

const DaysHeader = ({ }: Props) => {
    const ref = useSetCustomCssVariable({ '--scroll-margin-top': 'height' });
    const TranslatedPoolNames = [...translatedDays.values()] as [PoolName];

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
                    .
                </Box>
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(9, max-content)',
                        placeContent: 'center',
                        placeItems: 'center',
                    }}
                    className="links"
                >
                    {TranslatedPoolNames.map((day, i) => {
                        const activeClassColorSx: SxProps = {
                            borderColor: 'secondary.main',
                            '.active': {
                                backgroundColor: 'secondary.main',
                            },
                        };

                        return (
                            <Fragment key={day}>
                                {i === 0 && <Divider orientation="vertical" sx={activeClassColorSx} />}
                                <Box
                                    sx={{
                                        ...activeClassColorSx,
                                        transition: 'background-color 0.2s ease-in-out;',
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
                                <Divider orientation="vertical" sx={activeClassColorSx} />
                            </Fragment>
                        );
                    })}
                </Box>
            </Box>
        </>
    );
};

export default DaysHeader;
