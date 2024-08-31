import { Box, Button, Card, Paper, Typography } from "@mui/material";
import MyTickets from "../dashboard/MyTickets";
import MyEvents from "../dashboard/MyEvents";


export default function MyProfile() {
    return (
        <Paper>
        <Typography variant="h1">My Profile</Typography>
        <Typography variant="body1"> Hello world! There'll be stuff here at some point, but bear with me for now. Anyways, how's your day? I hope you're doing well.
             </Typography>
        <Typography variant="body1"> Anyways, here are the events you sent in. </Typography>
        <MyEvents />
        <Typography variant="body1"> And here are your tickets. I think? </Typography>
        <MyTickets />
        <Typography variant="body1"> Is there something wrong? </Typography>
        <Button variant="contained" color="primary" href="/">I'm not supposed to be here!</Button>
        </Paper>
    );
}