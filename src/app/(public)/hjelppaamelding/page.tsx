import BackButton from '$app/(authorized)/components/BackButton';
import { Card, CardContent, Container, List, ListItem, ListItemText, Paper, Typography } from '@mui/material';
import { Metadata } from 'next';
import Image from 'next/image';
import AwakeDragons from 'public/interessedragene/2024AwakeDragons1_1.png';
import HappyDragons from 'public/interessedragene/2024HappyDragons1_1.png';
import SleepyDragons from 'public/interessedragene/2024SleepyDragons1_1.png';
import VeryHappyDragons from 'public/interessedragene/2024VeryHappyDragons1_1.png';

export const metadata: Metadata = {
    title: 'Hjelp påmelding',
    description: 'Forklaring på hvordan pujepåmeldingen fungerer',
};

const HjelpPaameldingPage = async () => {
    return (
            <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
              <Typography variant="h3" gutterBottom>
                Slik fungerer påmeldingssystemet
              </Typography>
        
              <Typography variant="body1" paragraph>
                Det er viktig for oss at alle deltakarar på Regncon skal ha det kjekt på festival, og få spela i minst ei av puljene dei er mest interessert i. Puljepåmeldingssystemet vårt er meint å hjelpa til med dette.
              </Typography>
        
              <Card sx={{ mb: 4 }}>
                <CardContent>
                  <Typography variant="h5" gutterBottom>
                    Slik bruker du påmeldingssystemet:
                  </Typography>
                  <Typography variant="body1" paragraph>
                    Gå gjennom arrangementa, og vel om du er “litt interessert”, “interessert” eller “veldig interessert”, eventuelt “ikkje interessert” om den pulja ikkje er noko for deg.
                  </Typography>
                  <Typography variant="body1" paragraph>
                    Ver obs på at arrangement kan gå fleire gongar. Dersom du er veldig interessert i å delta på fredagen, men ikkje på laurdagen, kan du også velja dette. Er du veldig interessert uansett dag, vel du det.
                  </Typography>
                  <Typography variant="body1" paragraph>
                    Puljepåmeldinga blir stengt ei viss tid før puljene startar, slik at vi har tid til å setja saman puljene.
                  </Typography>
                </CardContent>
              </Card>
        
              <Card sx={{ mb: 4 }}>
                <CardContent>
                  <Typography variant="h5" gutterBottom>
                    Slik prioriterer vi:
                  </Typography>
                  <Typography variant="body1" paragraph>
                    Puljepåmeldingssystemet hjelper oss med å gjera puljefordelinga raskare, meir effektiv og meir rettvis. Samstundes vil skjønsvurderingar alltid spela ei rolle. Utover desse vert prioriteringane gjort slik:
                  </Typography>
        
                  <List>
                    <ListItem>
                        <Image src={VeryHappyDragons} alt="Veldig interessert" width={60} height={60} />
                      <ListItemText sx={{paddingLeft:"0.7rem"}} primary="1. Spelarar som ikkje tidlegare har fått delta i ei pulje dei er veldig interessert i, får tildelt plass." />
                    </ListItem>
                    <ListItem>
                        <Image src={VeryHappyDragons} alt="Veldig interessert" width={60} height={60} />
                      <ListItemText sx={{paddingLeft:"0.7rem"}} primary="2. Dersom det er ledige plassar, går desse til andre spelarar som har sagt at dei er veldig interessert." />
                    </ListItem>
                    <ListItem>
                        <Image src={HappyDragons} alt="Interessert" width={60} height={60} />
                      <ListItemText sx={{paddingLeft:"0.7rem"}} primary="3. Dersom det framleis er ledige plassar, går desse til spelarar som har sagt at dei er interessert." />
                    </ListItem>
                    <ListItem>
                        <Image src={AwakeDragons} alt="Litt interessert" width={60} height={60} />
                      <ListItemText sx={{paddingLeft:"0.7rem"}} primary="4. Er det endå ledige plassar, vert desse tildelt til spelarar som har sagt at dei er litt interessert." />
                    </ListItem>
                    <ListItem>
                      <Image src={SleepyDragons} alt="Ikke interessert" width={60} height={60} />
                      <ListItemText sx={{paddingLeft:"0.7rem"}}primary="5. Dersom nokon har sagt at dei ikkje er interessert, vert dei ikkje lagt til i pulja." />
                    </ListItem>
                    <ListItem>
                      <ListItemText primary="6. Er det framleis ledige plassar? Då vil desse vera tilgjengelege i samband med opprop." />
                    </ListItem>
                  </List>
                </CardContent>
              </Card>
            </Container>
    );
};
export default HjelpPaameldingPage;
