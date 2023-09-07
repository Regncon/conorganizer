import AppTopBar from "@/components/AppTopBar";
import AddEvent from "@/components/addEvent";
import { Box } from "@mui/material";
import { AuthProvider } from "@/components/auth";
import { Login } from "@mui/icons-material";
import EventList from "@/components/eventList";
import Menu from "@/components/menu";

export default function Home() {
  return (
    <main className="">
      <AppTopBar />

      <Box className="flex flex-row flex-wrap justify-center gap-4">
        <AuthProvider>
          <Login />
          <EventList />
        </AuthProvider>
      </Box>

    </main>
  );
}
