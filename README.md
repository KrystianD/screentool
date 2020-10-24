screentool
=====

A tool to make screenshots in a blazingly fast way.

# Demo

https://www.youtube.com/watch?v=1PFXvkRdBNw

# Basic usage

Use your favorite hotkey to start screentool

1. Select screenshot area

- Select a range by dragging mouse over area of interest - to make a screenshot of a range \
  OR
- Click a window - to make screenshot of a window \
  OR
- Click desktop or the edge of the screen - to make screenshot of the whole desktop

2. Release mouse button. Your screenshot is now saved to your clipboard.
   You can paste it somewhere, e.g. into Gimp, Hangouts or a Skype conversation.

## Advanced usage

Select screenshot range with one of above ways, press and hold Shift key and release left mouse button.
The tool will enter Annotation mode.

In Annotation mode, current tool can be changed with `Space`.

Currently, there are two tools implemented:

- Arrow - drag mouse to create an arrow
- Freehand drawing tool - drag mouse to create a freehand line

Release `Shift` to save the screenshot to the clipboard.

## Saving screenshot

screentool will save each captured screenshot in `$HOME/screenshots` directory if such directory is present
in the filesystem.

### Known limitations

Due to the nature of Go static linking, Go apps grow quickly in size.
As a result, if screentool is dropped out of filesystem cache, a noticeable lag may 
occur when starting the tool.