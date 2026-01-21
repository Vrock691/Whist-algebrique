import 'package:flutter/material.dart';
import 'package:whist_flutter_client/screens/home/JoinOrCreateGameScreen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Whist',
      theme: ThemeData(
        colorScheme: .fromSeed(seedColor: Colors.deepPurple),
      ),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

enum HomeTab {home, account}

class _HomePageState extends State<HomePage> {
  int tab = 0;

  void changeTab(int index) {
    setState(() {
      tab = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: (tab == 0) ? JoinCreateGameScreen() : SizedBox(),
      bottomNavigationBar: BottomNavigationBar(items: [
        BottomNavigationBarItem(icon: Icon(Icons.home), label: "Accueil"),
        BottomNavigationBarItem(icon: Icon(Icons.account_circle), label: "Profil"),
      ], onTap: (index) => changeTab(index), currentIndex: tab),
    );
  }
}
