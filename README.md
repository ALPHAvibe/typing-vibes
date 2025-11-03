# ⚡ Typing Vibes

A terminal-based typing speed test for Go developers. Practice typing by copying real Go functions from your own codebase with good vibes! 🎵

## Features

- 📁 **Use Your Own Code** - Point to any folder with Go files
- ⚡ **Real-time Stats** - Live WPM, accuracy, and progress tracking
- 🎯 **Smart Indentation** - Auto-skips leading whitespace when you press Enter
- 🎨 **Visual Feedback** - Dual-line display shows errors above the correct code
- ⚙️ **Customizable** - Configure function size, time limits, and folder paths
- 💾 **Persistent Config** - Settings saved to `~/.config/typing_vibes/`

## Installation
```bash
go install github.com/ALPHAvibe/typing-vibes@latest
```

Or clone and build:
```bash
git clone https://github.com/ALPHAvibe/typing-vibes.git
cd typing_vibes
go build
./typing_vibes
```

## Usage

1. **First run:** Press Enter to load a function from `~/code` (default)
2. **Start typing:** Type the function exactly as shown
3. **Press Enter:** Automatically skip indentation on new lines
4. **See your stats:** WPM, accuracy, and time tracking

## Keyboard Shortcuts

- `Enter` - Start/restart test
- `Ctrl+R` - Load new function
- `Ctrl+S` - Open settings
- `Esc` - Quit

## Configuration

Press `Ctrl+S` to configure:

- **Folder Path** - Where to find Go files
- **Min/Max Lines** - Function size range (default: 5-50)
- **Time Limit** - Max seconds per test (0 = unlimited, default: 30)

Config file: `~/.config/typing_vibes/typing_vibes.yaml`

## Screenshots

![Typing Vibes in action](https://via.placeholder.com/800x400.png?text=Add+your+screenshot+here)

## Why Typing Vibes?

- Practice typing with **real code patterns** from your projects
- Build **muscle memory** for common Go idioms
- Track your progress with **detailed statistics**
- Train on **your actual codebase**

