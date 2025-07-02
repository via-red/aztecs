import React from 'react';
import { Layout, Menu, theme } from 'antd';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'; // Import React Router components
import './App.css';

// Import page components
import BlockExplorer from './components/BlockExplorer';
import WalletManager from './components/WalletManager';
import TransactionSender from './components/TransactionSender';

const { Header, Content, Footer } = Layout;

function App() {
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  return (
    <Router> {/* Wrap with Router */}
      <Layout style={{ minHeight: '100vh' }}>
        <Header style={{ display: 'flex', alignItems: 'center' }}>
          <div className="demo-logo" />
          <Menu
            theme="dark"
            mode="horizontal"
            defaultSelectedKeys={['/']} // Use path as key
            items={[
              { key: '/', label: <Link to="/">区块链浏览器</Link> }, // Use Link for navigation
              { key: '/wallets', label: <Link to="/wallets">钱包管理</Link> },
              { key: '/send', label: <Link to="/send">发起交易</Link> },
            ]}
            style={{ flex: 1, minWidth: 0 }}
          />
        </Header>
        <Content style={{ padding: '0 48px' }}>
          <div
            style={{
              background: colorBgContainer,
              minHeight: 280,
              padding: 24,
              borderRadius: borderRadiusLG,
              marginTop: 16,
            }}
          >
            <Routes> {/* Define Routes */}
              <Route path="/" element={<BlockExplorer />} />
              <Route path="/wallets" element={<WalletManager />} />
              <Route path="/send" element={<TransactionSender />} />
              {/* Add a default route or 404 page if needed */}
            </Routes>
          </div>
        </Content>
        <Footer style={{ textAlign: 'center' }}>
          Simple Blockchain Project ©{new Date().getFullYear()} Created by Roo
        </Footer>
      </Layout>
    </Router>
  );
}

export default App;
