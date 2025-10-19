# Bingo Discord Bot

A Discord bot for hosting user-defined bingo games.

## Features

- Create custom bingo games from CSV event lists
- Distribute unique boards to players
- Vote on events as they occur
- Track game progress and winners

## Usage

1. Create a CSV file with your events (one per line)
2. Use `/new_game` to create a game from your CSV
3. Use `/set_active_game` to select which game to play
4. Players use `/view_board` to see their boards
5. Vote on events with `/vote` as they happen

## Tech Stack

- Go + [discordgo](https://github.com/bwmarrin/discordgo)
- SQLite database
