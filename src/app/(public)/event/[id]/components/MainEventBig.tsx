import type { ConEvent } from '$lib/types';
import { faUserSecret, faScroll } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { NavigateBefore } from '@mui/icons-material';
import { Paper, Box, Typography, Chip, Tooltip } from '@mui/material';
import MuiMarkdown from 'mui-markdown';
import Image from 'next/image';
import blekksprut2 from '$public/blekksprut2.jpg';
import InterestSelector from './components/InterestSelector';
import NavigatePreviousLink from './ui/NavigatePreviousLink';
import NavigateNextLink from './ui/NavigateNextLink';
import MuiMarkdownClient from './ui/MuiMarkdownClient';

type Props = {
    event: ConEvent;
    prevNavigationId?: string;
    nextNavigationId?: string;
};

const MainEventBig = ({ event, prevNavigationId, nextNavigationId }: Props) => {
    console.log(typeof window === 'undefined' ? 'server' : 'client');
    return (
        <Paper
            elevation={0}
            sx={{
                backgroundColor: 'black',
                '&': {
                    '--image-height': '15.1429rem',
                    '--slider-interest-width': '19.5714rem',
                    '--event-margin-left': '4rem',
                    '--event-header-margin-left': 'calc(var(--event-margin-left) - 1rem)',
                },
                maxWidth: '1200px',
            }}
        >
            <Box
                sx={{
                    display: 'grid',
                    '& > *': { gridColumn: '1 / 2', gridRow: '1 / 2' },
                }}
            >
                <Box
                    component={Image}
                    src={blekksprut2}
                    alt="noe alt-tekst"
                    sx={{
                        width: '100%',
                        height: 'auto',
                        maxHeight: 'var(--image-height)',
                        maxWidth: '100%',
                        // aspectRatio: '3.3 / 2',
                    }}
                    placeholder="blur"
                    loading="lazy"
                />
                <Box
                    sx={{
                        background: `linear-gradient(0deg, black, transparent)`,
                        maxHeight: 'var(--image-height)',
                    }}
                >
                    <Box
                        sx={{
                            display: 'grid',
                            height: '100%',
                            wordBreak: 'break-word',
                        }}
                    >
                        <Tooltip title={event.title || 'Tittel'}>
                            <Typography
                                variant="h1"
                                sx={{
                                    margin: '0',
                                    marginBlockStart: '1rem',
                                    marginInlineStart: 'var(--event-header-margin-left)',
                                    fontSize: 'clamp(1.7rem, 2.9vw, 3.42857rem)',
                                    whiteSpace: 'nowrap',
                                    textOverflow: 'ellipsis',
                                    overflow: 'clip',
                                    maxHeight: 'var(--image-height)',
                                    maxWidth:
                                        'min(calc(100dvw - var(--event-header-margin-left)), calc(1200px - var(--event-header-margin-left)))',
                                }}
                            >
                                {event.title || 'Tittel'}
                            </Typography>
                        </Tooltip>

                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: '50% 50%',
                                gridTemplateRows: '10.4286rem',
                            }}
                        >
                            <Box
                                sx={{
                                    display: 'grid',
                                    gridTemplateColum: '1fr 1fr',
                                }}
                            >
                                <Box
                                    sx={{
                                        display: 'grid',
                                        gridTemplateRows: '1fr 1fr',
                                        placeContent: 'start',
                                        marginInlineStart: 'var(--event-margin-left)',
                                    }}
                                >
                                    <Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
                                        <FontAwesomeIcon icon={faUserSecret} size="2x" style={{ color: '#ff7c7c' }} />
                                        <Box>
                                            <Typography component="span" sx={{ color: 'primary.main' }}>
                                                {event.icons?.includes('rollespill') ? 'Gamemaster' : 'Arrangør'}
                                            </Typography>
                                            <Typography variant="body1" margin={0}>
                                                {event.gameMaster || 'Navn'}
                                            </Typography>
                                        </Box>
                                    </Box>
                                    <Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
                                        <FontAwesomeIcon icon={faScroll} size="2x" style={{ color: '#ff7c7c' }} />
                                        <Box>
                                            <Typography component="span" sx={{ color: 'primary.main' }}>
                                                System
                                            </Typography>
                                            <Typography variant="body1" margin={0}>
                                                {event.system || 'System'}
                                            </Typography>
                                        </Box>
                                    </Box>
                                </Box>
                            </Box>

                            <Typography
                                sx={{
                                    overflow: 'clip',
                                    maxHeight: 'var(--image-height)',
                                }}
                            >
                                {event.shortDescription || 'Kort beskrivelse'}
                            </Typography>
                        </Box>
                    </Box>
                </Box>
            </Box>
            <Box sx={{ display: 'grid', gridAutoFlow: 'column', marginInline: '2rem', marginBlockStart: '0.5rem' }}>
                <Box sx={{ placeSelf: 'center start' }}>
                    <NavigatePreviousLink previousNavigationId={prevNavigationId} />
                </Box>
                <Box sx={{ placeSelf: 'center end' }}>
                    <NavigateNextLink nextNavigationId={nextNavigationId} />
                </Box>
            </Box>
            <Box
                sx={{
                    padding: '1.5rem 1rem 1rem 0rem',
                    display: 'grid',
                    gridTemplateColumns: '1fr 1fr',
                    marginInlineStart: 'var(--event-margin-left)',
                    position: 'relative',
                }}
            >
                <Box>
                    <Box
                        sx={{
                            display: 'flex',
                            flexWrap: 'wrap',
                            gap: '.5em',
                            overflowX: 'auto',
                            paddingBottom: '0.35rem',
                        }}
                    >
                        {/* {event.icons?.map((tag) => (
                                <Chip label={tag} key={tag} color="primary" icon={<NavigateBefore />} />
                            ))} */}

                        <Chip label={'tag'} key={'tag'} color="primary" variant="outlined" icon={<NavigateBefore />} />
                        <Chip
                            label={'tatagtagtagg'}
                            key={'tag'}
                            color="primary"
                            variant="outlined"
                            icon={<NavigateBefore />}
                        />
                        <Chip
                            label={'tagtagtag'}
                            key={'tag'}
                            color="primary"
                            variant="outlined"
                            icon={<NavigateBefore />}
                        />
                        <Chip
                            label={'tagtagtag'}
                            key={'tag'}
                            color="primary"
                            variant="outlined"
                            icon={<NavigateBefore />}
                        />
                        <Chip
                            label={'tagtagtag'}
                            key={'tag'}
                            color="primary"
                            variant="outlined"
                            icon={<NavigateBefore />}
                        />
                        <Chip
                            label={'tagtagtag'}
                            key={'tag'}
                            color="primary"
                            variant="outlined"
                            icon={<NavigateBefore />}
                        />
                    </Box>
                    <InterestSelector />
                </Box>
                <MuiMarkdownClient description={event.description} />
            </Box>
        </Paper>
    );
};

export default MainEventBig;
