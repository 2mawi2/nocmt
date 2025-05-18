import { useEffect, useState } from 'react';
import Notes from './notes/Notes';
import ShoppingList from './ShoppingList';
import { Gallery } from './gallery/Gallery';
import { FamilyTreeContainer } from './tree/FamilyTreeContainer';

interface MenuItem {
  name: string;
  icon: string;
  component: React.ReactNode;
}

interface DashboardProps {
  onLogout: () => void;
}

const menuItems: MenuItem[] = [
  { name: 'Family Tree', icon: '🌳', component: <FamilyTreeContainer /> },
  { name: 'Notes', icon: '📝', component: <Notes /> },
  { name: 'Shopping List', icon: '🛒', component: <ShoppingList /> },
  { name: 'Photos', icon: '📸', component: <Gallery /> },
];

export default function Dashboard({ onLogout }: DashboardProps) {
  const [logoAtTop, setLogoAtTop] = useState(false);
  const [showContent, setShowContent] = useState(false);
  const [activeMenuItem, setActiveMenuItem] = useState(menuItems[0].name);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  useEffect(() => {
    const hasAnimated = localStorage.getItem('hasAnimatedDashboard') === 'true';
    const justLoggedIn = localStorage.getItem('justLoggedIn') === 'true';

    if (justLoggedIn && !hasAnimated) {
      // Play animation only if just logged in and animation hasn't played
      const logoTimer = setTimeout(() => {
        setLogoAtTop(true);
      }, 1000);
      const contentTimer = setTimeout(() => {
        setShowContent(true);
        localStorage.setItem('hasAnimatedDashboard', 'true');
      }, 1700);
      localStorage.removeItem('justLoggedIn');
      return () => {
        clearTimeout(logoTimer);
        clearTimeout(contentTimer);
      };
    } else {
      // Skip animation and set states immediately
      setLogoAtTop(true);
      setShowContent(true);
    }
  }, []);

  const handleLogout = () => {
    document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT; secure; samesite=strict';
    onLogout();
  };

  return (
    <div className="relative w-screen h-screen bg-background overflow-hidden">
      {/* Logo */}
      <div
        className={`
          fixed transition-all duration-700 ease-in-out z-20
          ${logoAtTop
            ? 'top-4 left-4 md:left-4 text-2xl md:text-3xl'
            : 'top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-[clamp(2.5rem,10vw,5rem)]'
          }
          font-mono text-accent tracking-[0.15em]
          after:content-[''] after:inline-block after:align-bottom 
          after:w-[3px] after:h-[1em] after:bg-accent after:ml-[5px] 
          after:animate-cursor-blink
        `}
      >
        WICHTNER
      </div>

      {/* Mobile Menu Button */}
      <button
        onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
        className={`
          fixed top-4 right-4 z-30 p-2 text-accent md:hidden
          transition-opacity duration-500
          ${showContent ? 'opacity-100' : 'opacity-0'}
        `}
      >
        <svg 
          className="w-6 h-6" 
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          {isMobileMenuOpen ? (
            <path 
              strokeLinecap="round" 
              strokeLinejoin="round" 
              strokeWidth={2} 
              d="M6 18L18 6M6 6l12 12"
            />
          ) : (
            <path 
              strokeLinecap="round" 
              strokeLinejoin="round" 
              strokeWidth={2} 
              d="M4 6h16M4 12h16M4 18h16"
            />
          )}
        </svg>
      </button>

      {/* Sidebar */}
      <div
        className={`
          fixed top-0 left-0 h-full bg-background/95 backdrop-blur-sm
          border-r border-accent/20 transform transition-all duration-500 ease-in-out
          ${showContent ? 'translate-x-0' : '-translate-x-full'}
          md:w-64 md:translate-x-0
          ${isMobileMenuOpen ? 'w-full translate-x-0' : 'w-0 -translate-x-full'}
          pt-24 z-10 flex flex-col
        `}
      >
        <nav className="p-4 flex-grow overflow-y-auto">
          <ul className="space-y-4">
            {menuItems.map((item) => (
              <li key={item.name}>
                <button
                  onClick={() => {
                    setActiveMenuItem(item.name);
                    setIsMobileMenuOpen(false);
                  }}
                  className={`
                    w-full text-left px-4 py-3 rounded-lg transition-colors duration-200 
                    flex items-center gap-3
                    ${
                      activeMenuItem === item.name
                        ? 'bg-accent/20 text-accent'
                        : 'text-text-color hover:bg-accent/10'
                    }
                  `}
                >
                  <span className="text-xl">{item.icon}</span>
                  <span className="whitespace-nowrap">{item.name}</span>
                </button>
              </li>
            ))}
          </ul>
        </nav>

        {/* Logout Button */}
        <div className="p-4 border-t border-accent/20">
          <button 
            onClick={handleLogout}
            className="w-full px-4 py-3 text-text-color hover:bg-accent/10 rounded-lg transition-colors duration-200 flex items-center gap-3 hover:text-red-400"
          >
            <span className="text-xl">🚪</span>
            <span className="whitespace-nowrap">Logout</span>
          </button>
        </div>
      </div>

      {/* Overlay for mobile menu */}
      {isMobileMenuOpen && (
        <div 
          className="fixed inset-0 bg-black/50 z-0 md:hidden"
          onClick={() => setIsMobileMenuOpen(false)}
        />
      )}

      {/* Main content area */}
      <main
        className={`
          h-full transition-all duration-500
          ${showContent ? 'opacity-100' : 'opacity-0'}
          md:ml-64
          p-4 md:p-8 pt-20 md:pt-24
          relative
        `}
      >
        <div className="max-w-6xl mx-auto">
          {menuItems.find(item => item.name === activeMenuItem)?.component}
        </div>
      </main>

      {/* Scanline effect */}
      <div className="scanline"></div>
    </div>
  );
}