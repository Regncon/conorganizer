import type { PoolEvent } from '$lib/types';
import { faUserSecret, faScroll } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Paper, Box, Typography, Chip } from '@mui/material';
import InterestSelector from './components/InterestSelector';
import NavigatePreviousLink from './ui/NavigatePreviousLink';
import NavigateNextLink from './ui/NavigateNextLink';
import MuiMarkdownClient from './ui/MuiMarkdownClient';
import { createIconFromString } from '$app/(public)/components/lib/helpers/icons';
import GoToEventAdministrationButton from './ui/GoToEventAdministrationButton';
import type { InterestLevel } from '$lib/enums';

type Props = {
    poolEvent: PoolEvent;
    prevNavigationId?: string;
    nextNavigationId?: string;
    isAdmin?: boolean;
    activeParticipant?: { id?: string; interestLevel?: InterestLevel };
};

const MainEventBig = async ({
    poolEvent,
    prevNavigationId,
    nextNavigationId,
    isAdmin = false,
    activeParticipant,
}: Props) => {
    // console.log(typeof window === 'undefined' ? 'server' : 'client');

    return (
        <Paper
            elevation={0}
            sx={{
                backgroundColor: 'black',
                '&': {
                    '--image-height': '15.1429rem',
                    '--slider-interest-width': '23rem',
                    '--event-margin-left': '4rem',
                    '--event-header-margin-left': 'calc(var(--event-margin-left) - 1rem)',
                    '--event-width': 'min(100%, 1200px)',
                },
                maxWidth: 'var(--event-width)',
            }}
        >
            <Box
                sx={{
                    display: 'grid',
                    '& > *': { gridColumn: '1 / 2', gridRow: '1 / 2' },
                    '& > img': { width: ' var(--event-width)', height: '193px' },
                }}
            >
                <img
                    alt="Game logo"
                    src={poolEvent.bigImageURL ? poolEvent.bigImageURL : '/dice-big.webp'}
                    width={1200}
                    height={193}
                    sizes="100vw"
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
                        <Typography
                            variant="h1"
                            sx={{
                                margin: '0',
                                marginBlockStart: '1rem',
                                marginInlineStart: 'var(--event-header-margin-left)',
                                fontSize: 'clamp(1.7rem, 2.9vw, 3.42857rem)',
                                textShadow: '1px 0 0 #000, 0 -1px 0 #000, 0 1px 0 #000, -1px 0 0 #000',
                                whiteSpace: 'nowrap',
                                textOverflow: 'ellipsis',
                                overflow: 'clip',
                                maxHeight: 'var(--image-height)',
                                maxWidth:
                                    'min(calc(100dvw - var(--event-header-margin-left)), calc(1200px - var(--event-header-margin-left)))',
                            }}
                        >
                            {poolEvent.title || 'Tittel'}
                        </Typography>

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
                                                {poolEvent.gameType === 'rolePlaying' ? 'Gamemaster' : 'Arrangør'}
                                            </Typography>
                                            <Typography variant="body1" margin={0}>
                                                {poolEvent.gameMaster || 'Navn'}
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
                                                {poolEvent.system || 'System'}
                                            </Typography>
                                        </Box>
                                    </Box>
                                </Box>
                            </Box>

                            <Typography
                                sx={{
                                    overflow: 'clip',
                                    maxHeight: 'var(--image-height)',
                                    textShadow: '1px 0 0 #000, 0 -1px 0 #000, 0 1px 0 #000, -1px 0 0 #000',
                                }}
                            >
                                {poolEvent.shortDescription || 'Kort beskrivelse'}
                            </Typography>
                        </Box>
                    </Box>
                </Box>
            </Box>
            <Box sx={{ display: 'grid', gridAutoFlow: 'column', marginInline: '2rem', marginBlockStart: '0.5rem' }}>
                <Box sx={{ placeSelf: 'center start' }}>
                    <NavigatePreviousLink previousNavigationId={prevNavigationId} />
                </Box>

                <GoToEventAdministrationButton parentEventId={poolEvent.parentEventId} isAdmin={isAdmin} />

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
                        {poolEvent.icons?.map((iconOption) => (
                            <Chip
                                label={iconOption.label}
                                key={iconOption.label}
                                color="primary"
                                variant="outlined"
                                icon={createIconFromString(iconOption.iconName)}
                            />
                        ))}
                    </Box>
                    <InterestSelector
                        poolName={poolEvent.poolName}
                        poolEventId={poolEvent.id}
                        poolEventTitle={poolEvent.title}
                        activeParticipant={activeParticipant}
                    />
                </Box>
                <MuiMarkdownClient description={poolEvent.description} />
            </Box>
        </Paper>
    );
};

export default MainEventBig;
