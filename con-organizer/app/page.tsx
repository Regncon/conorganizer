import AppTopBar from "@/components/AppTopBar";
import EventUi from "@/components/eventUi";
import AddEvent from "@/components/addEvent";
import { Box } from "@mui/material";
import Experiment from "@/components/experiment";
import { AuthProvider } from "@/components/auth";
import Welcome from "@/components/welcome";
import { Login } from "@mui/icons-material";
import EventList from "@/components/eventList";

export default function Home() {
  return (
    <main className="">
      <AppTopBar />
      <Box className="flex flex-row flex-wrap justify-center gap-4">
        <AuthProvider>
          <Login />
          <Welcome />
          {/* <EventList /> */}
          <Experiment />
          <AddEvent />
        </AuthProvider>
      </Box>

    </main>
  );
}
