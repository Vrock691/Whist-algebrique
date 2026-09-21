import 'package:flutter/material.dart';

class JoinCreateGameScreen extends StatefulWidget {
  const JoinCreateGameScreen({super.key});

  @override
  State<JoinCreateGameScreen> createState() => _JoinCreateGameScreenState();
}

class _JoinCreateGameScreenState extends State<JoinCreateGameScreen> {
  final List<TextEditingController> _controllers = List.generate(
    4,
    (index) => TextEditingController(),
  );
  final List<FocusNode> _focusNodes = List.generate(4, (index) => FocusNode());

  @override
  void dispose() {
    for (var controller in _controllers) {
      controller.dispose();
    }
    for (var focusNode in _focusNodes) {
      focusNode.dispose();
    }
    super.dispose();
  }

  void _onTextChanged(int index, String value) {
    if (value.isNotEmpty && index < 3) {
      FocusScope.of(context).requestFocus(_focusNodes[index + 1]);
    } else if (value.isEmpty && index > 0) {
      FocusScope.of(context).requestFocus(_focusNodes[index - 1]);
    }
  }

  String _getEnteredCode() {
    return _controllers.map((controller) => controller.text).join();
  }

  void _validateAndJoin() {
    String code = _getEnteredCode();
    if (code.length == 4) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Rejoindre la partie avec le code: $code')),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Veuillez entrer un code de 4 caractères.'),
        ),
      );
    }
  }

  void _createGame() {}

  @override
  Widget build(BuildContext context) {
    return Padding(padding: EdgeInsetsGeometry.all(15), child: Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      spacing: 5,
      children: [
        Expanded(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text(
                'Entrez le code de la partie',
                style: TextStyle(fontSize: 20),
              ),
              const SizedBox(height: 20),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(4, (index) {
                  return Container(
                    width: 50,
                    height: 50,
                    margin: const EdgeInsets.symmetric(horizontal: 5),
                    child: TextField(
                      controller: _controllers[index],
                      focusNode: _focusNodes[index],
                      textAlign: TextAlign.center,
                      maxLength: 1,
                      keyboardType: TextInputType.text,
                      decoration: InputDecoration(
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(10),
                        ),
                        counterText: '',
                      ),
                      onChanged: (value) => _onTextChanged(index, value),
                    ),
                  );
                }),
              ),
            ],
          ),
        ),
        ElevatedButton(
          onPressed: _validateAndJoin,
          child: const Text('Rejoindre'),
        ),
        ElevatedButton(
          onPressed: _createGame,
          child: const Text("Créer une partie"),
        ),
      ],
    ));
  }
}
