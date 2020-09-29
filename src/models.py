import json
import logging

from watchdog.events import PatternMatchingEventHandler


# TODO: https://github.com/gorakhargosh/watchdog/blob/master/src/watchdog/events.py
class CsvFileWatcher(PatternMatchingEventHandler):
    def __init__(self, patterns="*.csv"):
        super(PatternMatchingEventHandler, self).__init__()

    def dispatch(self, event):
        super().dispatch(event)

    def on_created(self, event):
        logging.log(f"on_created event: {event}")

    def on_deleted(self, event):
        logging.log(f"on_deleted event: {event}")

    def on_modified(self, event):
        logging.log(f"on_cmodified event: {event}")

    def on_moved(self, event):
        logging.log(f"on_moved event: {event}")


class Name:
    def __init__(self, first_name, middle_name, last_name):
        self.first = first_name
        self.middle = middle_name
        self.last = last_name


class Person(object):
    def __init__(self, id=None, first_name=None, middle_name="", last_name=None, phone_number=None):
        self.id = id
        if id is None or len(id) == 0:
            raise ValueError("INTERNAL_ID Must be set")
        if first_name is None or len(first_name) == 0:
            raise ValueError("FIRST_NAME Must be set")
        if last_name is None or len(last_name) == 0:
            raise ValueError("LAST_NAME Must be set")
        if phone_number is None or len(phone_number) == 0:
            raise ValueError("PHONE_NUM Must be set")

        self.name = Name(first_name=first_name, middle_name=middle_name, last_name=last_name)
        self.phone = phone_number

    def toJSON(self) -> object:
        """
        toJSON will return an Object that can be written out to json with the
        json.dumps method. Remove's the middle name if it is not set.
        :return: object
        """
        name = {
            "first" : self.name.first,
            "middle": self.name.middle,
            "last"  : self.name.last,
        }

        if len(name["middle"]) == 0:
            del name["middle"]

        return {
            "id"   : self.id,
            "name" : name,
            "phone": self.phone
        }


class Error(object):
    def __init__(self, row, error):
        self.row = row
        self.message = str(error)

    def toJSON(self):
        return {
            'LINE_NUM' : self.row,
            'ERROR_MSG': self.message
        }
