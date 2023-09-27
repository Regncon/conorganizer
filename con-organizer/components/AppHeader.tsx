import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { Accordion, AccordionDetails, AccordionSummary, Box, Link } from '@mui/material';
import { Typography } from '@/lib/mui';

const AppHeader = () => {
    return (
        <Box sx={{ p: '2em', margin: '0 auto', maxWidth: '900px' }}>
            <header className="AppHeader">
                <img
                    src="/image/regnconlogony.png"
                    alt="Regncondragen for 2023"
                    className="regnconLogo"
                    onClick={() => (window.location.href = `/`)}
                />
                <div>
                    <Typography variant="h1" color="white">
                        Regncon XXXI
                    </Typography>
                    <Typography variant="h4">Program</Typography>
                    <Accordion>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panel1a-content"
                            id="panel1a-header"
                        >
                            <Typography>Kan du hjelpe oss?</Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Typography>
                                <span>
                                    Vi ha rekordstor påmelding i år og trenger flere spillledere.
                                </span>
                                <Link href="https://www.regncon.no/pamelding-av-arrangement/">Meld deg på her.</Link>
                            </Typography>
                        </AccordionDetails>
                    </Accordion>
                    {/*                     <Typography>
                        <span>Kan du hjelpe oss? Vi ha rekordstor påmelding i år og trenger flere spillledere.</span>
                        <Link href="https://www.regncon.no/pamelding-av-arrangement/">Meld deg på her.</Link>
                    </Typography> */}
                </div>
            </header>
        </Box>
    );
};

export default AppHeader;
