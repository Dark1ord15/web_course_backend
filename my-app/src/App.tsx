import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import MainPage from './pages/MainPage/MainPage.tsx';
import RoadPage from './pages/RoadPage/RoadPage.tsx';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MainPage />} />
        <Route path="/roads/:id" element={<RoadPage />} />
      </Routes>
    </Router>
  );
};

export default App;