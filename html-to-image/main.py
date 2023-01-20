from html2image import Html2Image
import sys


def main():
    if len(sys.argv) < 4:
        print("Usage: ", "./main.py", "input_html",
              "output_path", "width", "height")
        return
    hti = Html2Image(
        browser_executable='C:/Users/borna/Downloads/chromiun/chrome-win/chrome-win/chrome.exe',
        output_path=sys.argv[2],
        size=[int(sys.argv[3]), int(sys.argv[4])])
    print('opening chromium')
    print(sys.argv[1])
    hti.screenshot(html_file=sys.argv[1],
                   save_as='tmp.png')
    print('took screenshot successfully')


if __name__ == '__main__':
    main()
