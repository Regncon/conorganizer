import { Box, Button, Paper, TextField, Typography } from "@mui/material";
import Ticket from "./Ticket";
import CustomIcon from "$lib/icons/miscgame";

export default function MyTickets() {
    return (
        <Box sx={{display:"flex", paddingLeft:"2rem", gap:"1rem"}}>
            <Paper sx={{marginBottom: "2rem", paddingLeft:"2rem", maxWidth:"320px"}}>
                <Typography variant="h1">Bekreft e-post/Mine billetter</Typography>
                <Typography variant="h2">Har billetter</Typography>
                <Button variant="contained" color="primary">Bekreft e-post</Button>
                <Button variant="contained" color="primary">Log inn med Google</Button>
                <Typography variant="h2">Har ikke billetter</Typography>
                <Button variant="contained" color="primary">Kjøp billett</Button>
            </Paper>
            <Paper sx={{marginBottom: "2rem", paddingLeft:"2rem", maxWidth:"320px"}}>
                <Typography variant="h1">Ingen?/Mine billetter</Typography>
                <Typography variant="h2">Fant ingen billetter.</Typography>
                <Button variant="contained" color="primary">Kjøp billett</Button>
                <Button variant="contained" color="primary">Har allerede kjøpt billett</Button>
                <TextField label="Skriv inn e-posten du brukte på Checkin" />
                <Button variant="contained" color="primary">Hent billett</Button>
            </Paper>
            <Paper sx={{marginBottom: "2rem", paddingLeft:"2rem", width:"320px"}}>
                <Typography variant="h1">My Tickets</Typography>
                <CustomIcon color="primary" size="large" />
                <CustomIcon color="secondary" size="small" />
                <Ticket /> 
                <Ticket />
            </Paper>
        <Button variant="contained" color="primary" href="/my-profile">Go back to my profile</Button>
        </Box>
    );
}