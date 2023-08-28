import AppTopBar from "@/components/AppTopBar";
import EventUi from "@/components/eventUi";
import AddEvent from "@/components/addEvent";
import { Box } from "@mui/material";

export default function Home() {
  return (
    <main className="">
      <AppTopBar />
      <Box className="flex flex-row flex-wrap justify-center gap-4">
        <EventUi title="Title 1" />
        <EventUi title="Title 2" />
        <EventUi title="Title 3" />
      </Box>
      {/*       <div>
        <AuthProvider>
          <Login />
          <Welcome />
          <Experiment />
        </AuthProvider>
      </div> */}
      <AddEvent />
    </main>
  );
}
