import { faUserPlus } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
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
                </div>
                <Accordion sx={{ maxWidth: '20em', position: 'absolute', top: '0', right: '0', zIndex: '9000' }}>
                    <AccordionSummary
                        expandIcon={<ExpandMoreIcon />}
                        aria-controls="panel1a-content"
                        id="panel1a-header"
                    >
                        <Typography variant="caption">
                            <FontAwesomeIcon icon={faUserPlus} color="#55cc99" />
                            &nbsp; Regncon trenger deg!
                        </Typography>
                    </AccordionSummary>
                    <AccordionDetails>
                        <Typography>
                            Vi har rekordstor påmelding i år og trenger flere arrang&oslash;rer. Hvis du har noe du kan
                            arrangere selv,&nbsp;
                            <Link href="https://www.regncon.no/pamelding-av-arrangement/">meld det på her</Link>.
                            &nbsp;Hvis ikke, kan du se etter arrangementer med dette symbolet:&nbsp;
                            <FontAwesomeIcon icon={faUserPlus} color="#55cc99" /> og sende styret en epost hvis du vil
                            hjelpe med å arrangere: <Link href="mailto:arrangere@regncon.no">arrangere@regncon.no</Link>
                        </Typography>
                    </AccordionDetails>
                </Accordion>
            </header>
        </Box>
    );
};

export default AppHeader;
