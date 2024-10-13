import { Box, Divider, Typography } from '@mui/material';
import { getAllPoolEvents } from '../components/lib/serverAction';
import Logo from '../components/ui/Logo';
import { getTranslatedDay } from '../components/lib/helpers/translation';
import MainEventBig from '../event/[id]/components/MainEventBig';

export default async function Print() {
    const poolEvents = await getAllPoolEvents();

    return (
        <Box>
            <Logo />
            <Box>
                {[...poolEvents.entries()].map(([day, events]) => {
                    return (
                        <Box key={day}>
                            {events
                                .filter((e) => e.published === true)
                                .map((event) => (
                                    <>
                                        <Typography variant="h1" sx={{ color: 'black' }}>
                                            {getTranslatedDay(day)}
                                        </Typography>
                                        <MainEventBig key={event.id} poolEvent={event} />
                                        <Divider sx={{ pageBreakAfter: 'always' }} />
                                    </>
                                ))}
                        </Box>
                    );
                })}
            </Box>
        </Box>
    );
}
